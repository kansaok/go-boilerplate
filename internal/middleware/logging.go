package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/util"
	"github.com/kansaok/go-boilerplate/pkg/logger"
)

func LoggingMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()

        c.Next()

		duration := time.Since(start)
		logger.Info(c,"Request: ["+c.Request.Method+"] "+c.Request.URL.Path+" ("+c.ClientIP()+") | Status: "+strconv.Itoa(c.Writer.Status())+" | Latency: "+duration.String())
    }
}

func CustomRecovery(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		logger.ErrorLogger.Errorf("Panic recovered: %s", err)
		util.RespondWithError(c,util.CodeUnknownError,util.MESSAGES["ERROR"],err)
	}
	c.AbortWithStatus(http.StatusInternalServerError)
}
