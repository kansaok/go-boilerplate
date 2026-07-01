package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/controller"
)

func AuthRoutes(r *gin.RouterGroup) {
	api := r.Group("/auth")
    {
		api.POST("/register", controller.Register)
		api.POST("/login", controller.Login)
    }
}
