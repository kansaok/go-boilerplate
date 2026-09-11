package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BodyLimitMiddleware membungkus request body dengan http.MaxBytesReader
// untuk mencegah klien mengirim payload raksasa (DoS via memory exhaustion).
// Saat limit terlampaui, operasi baca body gagal dengan error
// *http.MaxBytesError yang bisa dideteksi handler via util.IsBodyTooLargeError.
func BodyLimitMiddleware(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		}
		c.Next()
	}
}