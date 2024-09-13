package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/config"
	"github.com/kansaok/go-boilerplate/internal/util"
)

func AuthMiddleware(jwtConfig *config.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenStr := c.GetHeader("Authorization")
        if tokenStr == "" {
            util.RespondWithError(c,util.CodeUnauthorized,util.MESSAGES["TOKEN_NOTFOUND"])
            c.Abort()
            return
        }

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtConfig.SecretKey), nil
        })

        if err != nil || !token.Valid {
            util.RespondWithError(c,util.CodeUnauthorized,util.MESSAGES["INVALID_TOKEN"])
            c.Abort()
            return
        }

        c.Next()
    }
}
