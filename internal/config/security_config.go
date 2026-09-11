package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// SecurityConfig stores security settings.
type SecurityConfig struct {
	AllowedHosts        []string
	SecureSSLRedirect   bool
	CSRFTokenSecure     bool
	SessionCookieSecure bool
	BrowserXSSFilter    bool
	ContentTypeNosniff  bool
}

// LoadSecurityConfigs loads security configurations from the .env file
func LoadSecurityConfigs() SecurityConfig {
	allowedHosts := strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")
	if len(allowedHosts) == 0 || (len(allowedHosts) == 1 && allowedHosts[0] == "") {
		allowedHosts = []string{"localhost"}
		log.Println("WARNING: ALLOWED_HOSTS not set, defaulting to localhost only")
	}

	secureSSLRedirect := parseBoolEnv("SECURE_SSL_REDIRECT", true)
	csrfTokenSecure := parseBoolEnv("CSRF_COOKIE_SECURE", true)
	sessionCookieSecure := parseBoolEnv("SESSION_COOKIE_SECURE", true)
	browserXSSFilter := parseBoolEnv("SECURE_BROWSER_XSS_FILTER", true)
	contentTypeNosniff := parseBoolEnv("SECURE_CONTENT_TYPE_NOSNIFF", true)

	return SecurityConfig{
		AllowedHosts:        allowedHosts,
		SecureSSLRedirect:   secureSSLRedirect,
		CSRFTokenSecure:     csrfTokenSecure,
		SessionCookieSecure: sessionCookieSecure,
		BrowserXSSFilter:    browserXSSFilter,
		ContentTypeNosniff:  contentTypeNosniff,
	}
}

func parseBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("WARNING: Could not parse %s as boolean, defaulting to %v", key, defaultValue)
		return defaultValue
	}
	return parsed
}
