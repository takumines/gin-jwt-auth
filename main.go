package main

import (
	"github.com/gin-gonic/gin"
	"github.com/takumines/gin-jwt-auth/db"
	"github.com/takumines/gin-jwt-auth/routes"
)

func main() {
	r := gin.Default()
	db.Init()
	routes.Setup(r)
	r.Run(":8080")
}
