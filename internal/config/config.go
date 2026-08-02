// Package config loads and validates application configuration from the
// environment.
package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// Config holds all runtime configuration. Sensitive fields are never logged
// or included in error messages.
type Config struct {
	Env            string        `env:"APP_ENV"`
	HTTPAddr       string        `env:"HTTP_ADDR"`
	DatabaseURL    string        `env:"DATABASE_URL" sensitive:"true"`
	SessionSecret  []byte        `env:"SESSION_SECRET" sensitive:"true"`
	OTPSecret      []byte        `env:"OTP_SECRET" sensitive:"true"`
	ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUTDOWN_TIMEOUT"`
	DevCodeVisible bool          `env:"APP_DEV_CODE_VISIBLE"`
	DemoMode       bool          `env:"DEMO_MODE"`
	LogLevel       string        `env:"LOG_LEVEL"`
}

// Load reads configuration from the environment and validates it. Error
// messages never reveal secret values.
func Load() (Config, error) {
	c := Config{
		Env:             getenv("APP_ENV", "development"),
		HTTPAddr:        getenv("HTTP_ADDR", "127.0.0.1:8080"),
		DatabaseURL:     getenv("DATABASE_URL", ""),
		ReadTimeout:     getduration("HTTP_READ_TIMEOUT", 15*time.Second),
		WriteTimeout:    getduration("HTTP_WRITE_TIMEOUT", 15*time.Second),
		IdleTimeout:     getduration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		ShutdownTimeout: getduration("HTTP_SHUTDOWN_TIMEOUT", 30*time.Second),
		DevCodeVisible:  getenv("APP_DEV_CODE_VISIBLE", "false") == "true",
		DemoMode:        getenv("DEMO_MODE", "false") == "true",
		LogLevel:        strings.ToLower(getenv("LOG_LEVEL", "info")),
	}
	c.SessionSecret = []byte(getenv("SESSION_SECRET", ""))
	c.OTPSecret = []byte(getenv("OTP_SECRET", ""))

	if err := c.validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

// IsProduction reports whether the environment is production.
func (c Config) IsProduction() bool { return c.Env == "production" }

func (c Config) validate() error {
	var errs []string
	if c.Env == "" {
		errs = append(errs, "APP_ENV is required")
	}
	switch c.Env {
	case "development", "staging", "production", "demo":
	default:
		errs = append(errs, fmt.Sprintf("APP_ENV %q is not one of development|staging|production|demo", c.Env))
	}

	if c.DemoMode {
		// Demo mode uses in-memory stores and fixed secrets. Skip DB
		// and secret validation so the demo can start with no env file.
		if c.DevCodeVisible && c.IsProduction() {
			errs = append(errs, "APP_DEV_CODE_VISIBLE must be false in production")
		}
		if len(errs) > 0 {
			return errors.New("config: " + strings.Join(errs, "; "))
		}
		return nil
	}

	if c.DatabaseURL == "" {
		errs = append(errs, "DATABASE_URL is required")
	} else if _, err := url.Parse(c.DatabaseURL); err != nil {
		errs = append(errs, "DATABASE_URL is not a valid URL")
	}
	if len(c.SessionSecret) < 32 {
		errs = append(errs, "SESSION_SECRET must be at least 32 bytes")
	}
	if len(c.OTPSecret) < 32 {
		errs = append(errs, "OTP_SECRET must be at least 32 bytes")
	}
	if c.DevCodeVisible && c.IsProduction() {
		errs = append(errs, "APP_DEV_CODE_VISIBLE must be false in production")
	}
	if len(errs) > 0 {
		return errors.New("config: " + strings.Join(errs, "; "))
	}
	return nil
}

func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getduration(key string, fallback time.Duration) time.Duration {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return fallback
}