package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/controllers"
)

func Setup(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("login", controllers.Login)
	r.GET("/user", controllers.User)
	r.GET("/logout", controllers.Logout)
	r.POST("/forget", controllers.Forget)
	r.POST("/reset", controllers.Reset)
}
