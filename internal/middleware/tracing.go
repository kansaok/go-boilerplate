package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
)

var excludedRoutes = map[string]bool{
    "/product": true,
    "/marl":    true,
    "/inces":   true,
}

func TracingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
		// Check if the current route is excluded
        if _, ok := excludedRoutes[c.FullPath()]; ok {
            c.Next() // Skip middleware for excluded routes
            return
        }
        tracer := otel.Tracer("myapp/server")

        ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
        defer span.End()

        // Tambahkan context dengan span
        c.Request = c.Request.WithContext(ctx)

        c.Next()
    }
}
