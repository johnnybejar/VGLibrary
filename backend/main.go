package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {fmt.Println("GET path from root")})

	router.Run("localhost:8000");
}