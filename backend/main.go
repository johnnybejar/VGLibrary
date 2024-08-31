package main

import (
	"backend/controllers"
	"backend/initializers"
	"backend/middleware"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitDBConn()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {fmt.Println("GET path from root")})
	router.GET("/access-token", middleware.RequireAuth, controllers.GetAccessToken)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.GET("/logout", middleware.RequireAuth, controllers.Logout)
	router.GET("/validate", middleware.RequireAuth, controllers.Validate)
	
	router.SetTrustedProxies(nil)
	router.Run("localhost:8000");
}