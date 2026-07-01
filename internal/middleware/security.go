package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net"
	"net/http"
	"strings"

	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/util"

	"github.com/gin-gonic/gin"
)

func isHostAllowed(host string) bool {
	allowedHosts := config.LoadConfig().SecurityConfig.AllowedHosts
	hostWithoutPort := getHostWithoutPort(host)

	for _, allowedHost := range allowedHosts {
		if hostWithoutPort == allowedHost {
			return true
		}
	}
	return false
}

func getHostWithoutPort(host string) string {
	if strings.Contains(host, ":") {
		host, _, _ = net.SplitHostPort(host)
	}
	return host
}

func ValidateHost(c *gin.Context) {
	host := c.Request.Host
	if !isHostAllowed(host) {
		util.RespondWithError(c, util.CodeForbidden, util.MESSAGES["HOST_NOT_ALLOWED"], nil)
		c.Abort()
		return
	}
	c.Next()
}

func EnforceSSLRedirect(c *gin.Context) {
	if config.LoadConfig().SecurityConfig.SecureSSLRedirect {
		if c.Request.TLS == nil {
			target := "https://" + c.Request.Host + c.Request.RequestURI
			c.Redirect(http.StatusMovedPermanently, target)
			c.Abort()
			return
		}
	}
	c.Next()
}

func SetCSRFHeaders(c *gin.Context) {
	if config.LoadConfig().SecurityConfig.CSRFTokenSecure {
		csrfToken := generateCSRFToken()
		c.Header("Set-Cookie", "csrf_token="+csrfToken+"; Secure; HttpOnly; SameSite=Strict")
		c.Set("csrf_token", csrfToken)
	}
	c.Next()
}

func SetSessionCookie(c *gin.Context) {
	sessionID := generateSessionID()
	if config.LoadConfig().SecurityConfig.SessionCookieSecure {
		c.SetCookie("session_id", sessionID, 3600, "/", "", true, true)
	} else {
		c.SetCookie("session_id", sessionID, 3600, "/", "", false, false)
	}
	c.Set("session_id", sessionID)
	c.Next()
}

func SetSecurityHeaders(c *gin.Context) {
	// Content Security Policy
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; frame-ancestors 'none';")

	// Strict Transport Security
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

	// Referrer Policy
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	// Permissions Policy
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

	// X-Frame-Options
	c.Header("X-Frame-Options", "DENY")

	c.Next()
}

func SetXSSFilterHeader(c *gin.Context) {
	if config.LoadConfig().SecurityConfig.BrowserXSSFilter {
		c.Header("X-XSS-Protection", "1; mode=block")
	}
	c.Next()
}

func SetContentTypeNosniffHeader(c *gin.Context) {
	if config.LoadConfig().SecurityConfig.ContentTypeNosniff {
		c.Header("X-Content-Type-Options", "nosniff")
	}
	c.Next()
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func generateSessionID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
