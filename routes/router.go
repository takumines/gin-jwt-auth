package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/controllers"
)

func Setup(r *gin.Engine) {
	r.GET("/", controllers.Home)
}
