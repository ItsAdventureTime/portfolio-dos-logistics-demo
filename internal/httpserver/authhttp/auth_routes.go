// Package authhttp wires the auth HTTP routes onto the Chi router. It
// handles login, email verification, logout, session info, role preview,
// and the current-user endpoint. All error responses are neutral and do
// not reveal whether an account exists.
package authhttp

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/httpserver/middleware"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/service"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	svc    *service.AuthService
	otpCfg auth.OTPConfig
}

// Mount registers the auth routes on the given router.
func Mount(r chi.Router, svc *service.AuthService, otpCfg auth.OTPConfig) {
	h := &handler{svc: svc, otpCfg: otpCfg}
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/verify-email", h.VerifyEmail)
		r.Group(func(r chi.Router) {
			r.Use(middleware.RequireAuth(svc.ValidateSession))
			r.Use(middleware.CSRF)
			r.Post("/logout", h.Logout)
			r.Get("/session", h.Session)
			r.Post("/role-preview", h.SetRolePreview)
		})
	})
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth(svc.ValidateSession))
		r.Get("/auth/me", h.Me)
	})
}

type loginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type verifyEmailRequest struct {
	Identifier string `json:"identifier"`
	Code       string `json:"code"`
}

type rolePreviewRequest struct {
	RolePreview string `json:"role_preview"`
}

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := errorResponse{}
	resp.Error.Code = code
	resp.Error.Message = message
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "The request body could not be read.")
		return
	}
	if req.Identifier == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "Username or email and password are required.")
		return
	}
	ip := clientIP(r)
	userAgent := r.UserAgent()

	result, err := h.svc.Login(r.Context(), req.Identifier, req.Password, ip, userAgent)
	if err != nil {
		// All auth failures return the same neutral message.
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeError(w, http.StatusUnauthorized, "invalid_credentials",
				"The credentials provided do not match an active account.")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error",
			"An error occurred while processing the request.")
		return
	}

	if result.NeedsOTP {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":        "otp_required",
			"display_name":  result.DisplayName,
			"message":       "A verification code has been sent to the email on file.",
		})
		return
	}

	csrfToken, err := auth.GenerateCSRFToken(auth.HashToken(result.SessionToken))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "An error occurred.")
		return
	}
	middleware.SetSessionCookie(w, result.SessionToken, csrfToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":       "authenticated",
		"display_name": result.DisplayName,
		"csrf_token":   csrfToken,
	})
}

func (h *handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req verifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "The request body could not be read.")
		return
	}
	if req.Identifier == "" || req.Code == "" {
		writeError(w, http.StatusBadRequest, "invalid_request", "Identifier and verification code are required.")
		return
	}
	ip := clientIP(r)
	userAgent := r.UserAgent()

	result, err := h.svc.VerifyEmail(r.Context(), req.Identifier, req.Code, ip, userAgent)
	if err != nil {
		if errors.Is(err, service.ErrOTPInvalid) {
			writeError(w, http.StatusUnauthorized, "otp_invalid",
				"The verification code is invalid or has expired. Please request a new code.")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "An error occurred.")
		return
	}

	csrfToken, err := auth.GenerateCSRFToken(auth.HashToken(result.SessionToken))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "An error occurred.")
		return
	}
	middleware.SetSessionCookie(w, result.SessionToken, csrfToken)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":       "authenticated",
		"display_name": result.DisplayName,
		"csrf_token":   csrfToken,
	})
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(middleware.SessionCookieNameFor())
	if err == nil && cookie.Value != "" {
		tokenHash := auth.HashToken(cookie.Value)
		_ = h.svc.Logout(r.Context(), tokenHash)
	}
	middleware.ClearSessionCookie(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]any{"status": "logged_out"})
}

func (h *handler) Session(w http.ResponseWriter, r *http.Request) {
	ac, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id":      string(ac.User.ID),
		"username":     ac.User.Username,
		"display_name": ac.User.DisplayName,
		"email":        ac.User.Email,
		"role_preview": string(ac.Session.RolePreview),
	})
}

func (h *handler) SetRolePreview(w http.ResponseWriter, r *http.Request) {
	ac, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required.")
		return
	}
	var req rolePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", "The request body could not be read.")
		return
	}
	preview := domain.RolePreview(req.RolePreview)
	if !domain.IsValidRolePreview(preview) {
		writeError(w, http.StatusBadRequest, "invalid_role_preview",
			"The selected role preview is not valid.")
		return
	}
	tokenHash := auth.HashToken(ac.Session.TokenHash)
	if err := h.svc.SetRolePreview(r.Context(), tokenHash, preview); err != nil {
		if errors.Is(err, service.ErrSessionInvalid) {
			writeError(w, http.StatusUnauthorized, "session_expired",
				"The session has expired. Please sign in again.")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "An error occurred.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":       "role_preview_set",
		"role_preview": string(preview),
	})
}

func (h *handler) Me(w http.ResponseWriter, r *http.Request) {
	ac, ok := middleware.AuthFromContext(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Authentication required.")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"user_id":       string(ac.User.ID),
		"username":      ac.User.Username,
		"display_name":  ac.User.DisplayName,
		"email":         ac.User.Email,
		"email_verified": ac.User.EmailVerified,
		"role_preview":  string(ac.Session.RolePreview),
	})
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.TrimSpace(strings.Split(xff, ",")[0])
	}
	remote := r.RemoteAddr
	if idx := strings.LastIndex(remote, ":"); idx != -1 {
		remote = remote[:idx]
	}
	return remote
}