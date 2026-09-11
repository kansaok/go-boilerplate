package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/util"
)

func TestBodyLimitRejectsOversizedPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BodyLimitMiddleware(4))
	r.POST("/test", func(c *gin.Context) {
		if _, err := io.ReadAll(c.Request.Body); util.IsBodyTooLargeError(err) {
			c.Status(http.StatusRequestEntityTooLarge)
			return
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("abcdefghij"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected 413, got %d", w.Code)
	}
}

func TestBodyLimitAllowsWithinLimit(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(BodyLimitMiddleware(4096))
	r.POST("/test", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil || len(body) != 2 {
			c.Status(http.StatusBadRequest)
			return
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader("ok"))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestIsBodyTooLargeError(t *testing.T) {
	if util.IsBodyTooLargeError(nil) {
		t.Fatal("nil bukan error body too large")
	}
	if util.IsBodyTooLargeError(http.ErrBodyNotAllowed) {
		t.Fatal("error lain tidak boleh terdeteksi sebagai body too large")
	}
	if !util.IsBodyTooLargeError(&http.MaxBytesError{Limit: 16}) {
		t.Fatal("MaxBytesError harus terdeteksi")
	}
}