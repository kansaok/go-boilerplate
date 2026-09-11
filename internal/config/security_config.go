package config

import (
	"log"
	"net"
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
	MaxBodyBytes        int64
	MetricsUser         string
	MetricsPassword     string
	MetricsAllowlist    []*net.IPNet
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

	maxBodyBytes := int64(16 << 20)
	if raw := os.Getenv("MAX_BODY_BYTES"); raw != "" {
		if parsed, err := strconv.ParseInt(raw, 10, 64); err != nil || parsed <= 0 {
			log.Printf("WARNING: Could not parse MAX_BODY_BYTES (%q), defaulting to 16MB", raw)
		} else {
			maxBodyBytes = parsed
		}
	}

	return SecurityConfig{
		AllowedHosts:        allowedHosts,
		SecureSSLRedirect:   secureSSLRedirect,
		CSRFTokenSecure:     csrfTokenSecure,
		SessionCookieSecure: sessionCookieSecure,
		BrowserXSSFilter:    browserXSSFilter,
		ContentTypeNosniff:  contentTypeNosniff,
		MaxBodyBytes:        maxBodyBytes,
		MetricsUser:         os.Getenv("METRICS_USER"),
		MetricsPassword:     os.Getenv("METRICS_PASSWORD"),
		MetricsAllowlist:    parseCIDRList(os.Getenv("METRICS_ALLOWLIST")),
	}
}

// parseCIDRList parses a comma-separated list of CIDRs. Invalid entries are
// skipped with a warning so one typo does not disable the whole list.
func parseCIDRList(raw string) []*net.IPNet {
	networks := make([]*net.IPNet, 0)
	for _, part := range strings.Split(raw, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		_, ipNet, err := net.ParseCIDR(part)
		if err != nil {
			log.Printf("WARNING: invalid CIDR %q in METRICS_ALLOWLIST, skipped", part)
			continue
		}
		networks = append(networks, ipNet)
	}
	return networks
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
