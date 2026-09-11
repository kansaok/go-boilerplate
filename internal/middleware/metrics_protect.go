package middleware

import (
	"crypto/subtle"
	"net"

	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/util"
)

// MetricsProtection membatasi akses ke /metrics, dalam urutan prioritas:
//  1. Jika METRICS_USER/METRICS_PASSWORD dikonfigurasi -> wajib HTTP Basic Auth.
//  2. Selain itu, hanya mengizinkan IP loopback / CIDR pada METRICS_ALLOWLIST.
//
// Default (tanpa konfigurasi) menolak semua akses dari luar loopback sehingga
// inventori route dan volume trafik tidak bocor ke publik.
func MetricsProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		sec := config.LoadConfig().SecurityConfig

		if sec.MetricsUser != "" {
			user, pass, ok := c.Request.BasicAuth()
			if !ok ||
				subtle.ConstantTimeCompare([]byte(user), []byte(sec.MetricsUser)) != 1 ||
				subtle.ConstantTimeCompare([]byte(pass), []byte(sec.MetricsPassword)) != 1 {
				c.Header("WWW-Authenticate", `Basic realm="metrics"`)
				util.RespondWithError(c, util.CodeUnauthorized, "metrics authentication required", nil)
				c.Abort()
				return
			}
			c.Next()
			return
		}

		ip := net.ParseIP(c.ClientIP())
		if ip == nil {
			util.RespondWithError(c, util.CodeForbidden, "metrics endpoint not accessible from this network", nil)
			c.Abort()
			return
		}

		allowed := ip.IsLoopback()
		if !allowed {
			for _, cidr := range sec.MetricsAllowlist {
				if cidr.Contains(ip) {
					allowed = true
					break
				}
			}
		}

		if !allowed {
			util.RespondWithError(c, util.CodeForbidden, "metrics endpoint not accessible from this network", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}