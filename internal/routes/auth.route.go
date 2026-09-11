package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kansaok/go-boilerplate/internal/controller"
)

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
}
