package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// SecurityConfig stores security settings.
type SecurityConfig struct {
	AllowedHosts          []string
	SecureSSLRedirect     bool
	CSRFTokenSecure       bool
	SessionCookieSecure   bool
	BrowserXSSFilter      bool
	ContentTypeNosniff    bool
}

// LoadSecurityConfigs loads security configurations from the .env file
func LoadSecurityConfigs() SecurityConfig {
	allowedHosts := strings.Split(os.Getenv("ALLOWED_HOSTS"), ",")
	secureSSLRedirect, err := strconv.ParseBool(os.Getenv("SECURE_SSL_REDIRECT"))
	if err != nil {
		log.Printf("Error parsing SECURE_SSL_REDIRECT: %v, defaulting to false", err)
		secureSSLRedirect = false
	}

	csrfTokenSecure, err := strconv.ParseBool(os.Getenv("CSRF_COOKIE_SECURE"))
	if err != nil {
		log.Printf("Error parsing CSRF_COOKIE_SECURE: %v, defaulting to false", err)
		csrfTokenSecure = false
	}

	sessionCookieSecure, err := strconv.ParseBool(os.Getenv("SESSION_COOKIE_SECURE"))
	if err != nil {
		log.Printf("Error parsing SESSION_COOKIE_SECURE: %v, defaulting to false", err)
		sessionCookieSecure = false
	}

	browserXSSFilter, err := strconv.ParseBool(os.Getenv("SECURE_BROWSER_XSS_FILTER"))
	if err != nil {
		log.Printf("Error parsing SECURE_BROWSER_XSS_FILTER: %v, defaulting to false", err)
		browserXSSFilter = false
	}

	contentTypeNosniff, err := strconv.ParseBool(os.Getenv("SECURE_CONTENT_TYPE_NOSNIFF"))
	if err != nil {
		log.Printf("Error parsing SECURE_CONTENT_TYPE_NOSNIFF: %v, defaulting to false", err)
		contentTypeNosniff = false
	}

	return SecurityConfig{
		AllowedHosts:        allowedHosts,
		SecureSSLRedirect:   secureSSLRedirect,
		CSRFTokenSecure:     csrfTokenSecure,
		SessionCookieSecure: sessionCookieSecure,
		BrowserXSSFilter:    browserXSSFilter,
		ContentTypeNosniff:  contentTypeNosniff,
	}
}
