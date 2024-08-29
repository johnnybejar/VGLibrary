package main

import (
	"backend/controllers"
	"backend/initializers"
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.InitDBConn()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {fmt.Println("GET path from root")})
	router.GET("/access-token", getAccessToken)
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.GET("/logout", controllers.Logout)
	
	router.SetTrustedProxies(nil)
	router.Run("localhost:8000");
}

func getAccessToken(c *gin.Context) {
	godotenv.Load("../.env")
	CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("TWITCH_CLIENT_SECRET")
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", CLIENT_ID, CLIENT_SECRET)
	
	postBody, err := json.Marshal(&models.AccessToken{})
	if err != nil {
		log.Fatal(err);
	}
	
	reader := bytes.NewReader(postBody)

	res, err := http.Post(url, "", reader)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	if res.StatusCode >= 400 {
		c.JSON(res.StatusCode, gin.H{"Error response, Status Code": res.StatusCode})
		return
	}

	var body models.AccessToken
	
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, body)
}