package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
)

// initMetricsEnv menyiapkan env minimum agar config.LoadConfig() tidak fatal
// saat MetricsProtection dijalankan pada test (config di-cache per proses).
func initMetricsEnv() {
	os.Setenv("DB_CONNECTION", "sqlite")
	os.Setenv("DB_FILE", "/tmp/metrics-test.db")
	os.Setenv("JWT_SECRET_KEY", strings.Repeat("s", 40))
	os.Setenv("ALLOWED_HOSTS", "localhost")
	os.Setenv("METRICS_USER", "")
	os.Setenv("METRICS_PASSWORD", "")
	os.Setenv("METRICS_ALLOWLIST", "")
}

func setupMetricsRouter() *gin.Engine {
	config.ResetConfigTestOnly()
	gin.SetMode(gin.TestMode)
	r := gin.New()
	_ = r.SetTrustedProxies(nil)
	r.GET("/metrics", MetricsProtection(), func(c *gin.Context) { c.Status(http.StatusOK) })
	return r
}

func TestMetricsProtectionAllowsLoopbackByDefault(t *testing.T) {
	initMetricsEnv()
	r := setupMetricsRouter()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.RemoteAddr = "127.0.0.1:9999"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("loopback harus diizinkan secara default, got %d", w.Code)
	}
}

func TestMetricsProtectionDeniesRemoteByDefault(t *testing.T) {
	initMetricsEnv()
	r := setupMetricsRouter()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	req.RemoteAddr = "198.51.100.7:9999"
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("IP remote harus ditolak secara default, got %d", w.Code)
	}
}

func TestMetricsProtectionRequiresBasicAuth(t *testing.T) {
	initMetricsEnv()
	os.Setenv("METRICS_USER", "prom")
	os.Setenv("METRICS_PASSWORD", "s3cr3t")
	r := setupMetricsRouter()

	t.Run("tanpa kredensial", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		req.RemoteAddr = "127.0.0.1:9999"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 tanpa kredensial, got %d", w.Code)
		}
	})

	t.Run("kredensial benar", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		req.RemoteAddr = "198.51.100.7:9999"
		req.SetBasicAuth("prom", "s3cr3t")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("expected 200 dengan kredensial benar, got %d", w.Code)
		}
	})

	t.Run("kredensial salah", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
		req.RemoteAddr = "127.0.0.1:9999"
		req.SetBasicAuth("prom", "wrong")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Fatalf("expected 401 dengan kredensial salah, got %d", w.Code)
		}
	})
}