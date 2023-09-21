package main

import (
	"github.com/gimtwi/go-jwt-auth/controllers"
	"github.com/gimtwi/go-jwt-auth/initializers"
	"github.com/gimtwi/go-jwt-auth/middleware"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
	initializers.MigrateDB()
}

func main() {
	r := gin.Default()

	r.POST("/signup", controllers.SignUp)
	r.POST("/login", controllers.Login)
	r.GET("/validate", middleware.RequireAuth, controllers.Validate)

	r.Run()
}
