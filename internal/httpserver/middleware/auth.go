package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/auth"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/domain"
	"github.com/ItsAdventureTime/portfolio-dos-logistics-demo/internal/observability"
)

type authCtxKey struct{}

// SessionCookieName is the name of the session cookie.
const SessionCookieName = "__Host-dos_session"

// CSRFCookieName is the name of the CSRF double-submit cookie.
const CSRFCookieName = "__Host-dos_csrf"

// CSRFHeaderName is the request header that must carry the CSRF token.
const CSRFHeaderName = "X-CSRF-Token"

// AuthContext holds the authenticated session and user for the current
// request. Populated by RequireAuth and available to handlers via
// AuthFromContext.
type AuthContext struct {
	Session *domain.Session
	User    *domain.User
}

// AuthFromContext extracts the AuthContext from the request context.
func AuthFromContext(ctx context.Context) (*AuthContext, bool) {
	ac, ok := ctx.Value(authCtxKey{}).(*AuthContext)
	return ac, ok
}

// RequireAuth wraps the auth service to validate the session cookie on
// every request. If the session is invalid or absent, it returns 401 for
// API clients or redirects to login for browser requests. It populates
// the request context with AuthContext on success.
func RequireAuth(validate func(ctx context.Context, tokenHash string) (*domain.Session, *domain.User, error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil || cookie.Value == "" {
				deny(w, r)
				return
			}
			tokenHash := auth.HashToken(cookie.Value)
			sess, user, err := validate(r.Context(), tokenHash)
			if err != nil || sess == nil || user == nil {
				clearSessionCookie(w)
				deny(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), authCtxKey{}, &AuthContext{
				Session: sess,
				User:    user,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// CSRF protects state-changing methods (POST, PUT, PATCH, DELETE) via the
// double-submit cookie pattern. The CSRF token cookie is set on successful
// auth and must be echoed in the X-CSRF-Token header. GET/HEAD/OPTIONS are
// exempt.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isStateChanging(r.Method) {
			cookie, err := r.Cookie(CSRFCookieName)
			if err != nil || cookie.Value == "" {
				http.Error(w, `{"error":{"code":"csrf_failed","message":"CSRF token missing"}}`, http.StatusForbidden)
				return
			}
			headerToken := r.Header.Get(CSRFHeaderName)
			if headerToken == "" {
				http.Error(w, `{"error":{"code":"csrf_failed","message":"CSRF token missing"}}`, http.StatusForbidden)
				return
			}
			if !auth.VerifyCSRFToken(headerToken, cookie.Value) {
				http.Error(w, `{"error":{"code":"csrf_failed","message":"CSRF token invalid"}}`, http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func isStateChanging(method string) bool {
	switch strings.ToUpper(method) {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

func deny(w http.ResponseWriter, r *http.Request) {
	if isBrowserRequest(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_, _ = w.Write([]byte(`{"error":{"code":"unauthorized","message":"Authentication required"}}`))
}

func isBrowserRequest(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	return strings.Contains(accept, "text/html")
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// SetSessionCookie sets the session cookie and a CSRF double-submit cookie.
func SetSessionCookie(w http.ResponseWriter, sessionToken, csrfToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    csrfToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// ClearSessionCookie removes the session and CSRF cookies on logout.
func ClearSessionCookie(w http.ResponseWriter) {
	clearSessionCookie(w)
	http.SetCookie(w, &http.Cookie{
		Name:     CSRFCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

// _ retains the observability import for future correlation in middleware.
var _ = observability.CorrelationFrom