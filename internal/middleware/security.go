package middleware

import (
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
		util.RespondWithError(c,util.CodeForbidden,util.MESSAGES["HOST_NOT_ALLOWED"],nil)

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
		c.Header("Set-Cookie", "csrf_token=your-token; Secure; HttpOnly; SameSite=Strict")
	}
	c.Next()
}

func SetSessionCookie(c *gin.Context) {
	if config.LoadConfig().SecurityConfig.SessionCookieSecure {
		c.SetCookie("session_id", "your-session-id", 3600, "/", "", true, true)
	} else {
		c.SetCookie("session_id", "your-session-id", 3600, "/", "", false, false)
	}
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
