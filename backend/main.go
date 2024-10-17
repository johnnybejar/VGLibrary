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

	// Auth Routes
	router.GET("/", func(ctx *gin.Context) {fmt.Println("GET path from root")})
	auth := router.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/logout", middleware.RequireAuth, controllers.Logout)
		auth.GET("/validate", middleware.RequireAuth, controllers.Validate)
	}

	// API Routes
	giantbomb := router.Group("/giantbomb")
	{
		giantbomb.GET("/game", middleware.RequireAuth, controllers.GetGame)
		giantbomb.GET("/search", middleware.RequireAuth, controllers.SearchGames)
	}
	
	// DB Routes
	db := router.Group("/db")
	{
		db.POST("/create-game", controllers.WriteGame)
		db.GET("/game")
	}
	
	router.SetTrustedProxies(nil)
	router.Run("localhost:8000")
}