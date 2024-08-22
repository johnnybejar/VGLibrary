package main

import (
	. "backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", func(ctx *gin.Context) {fmt.Println("GET path from root")})
	router.GET("/access-token", getAccessToken).Use(AuthRequired())
	router.POST("/login", login)
	router.GET("/logout", logout)
	
	priv := router.Group("/private")
	{
		priv.GET("/", private)
		priv.GET("/status", status)
	}
	priv.Use(AuthRequired())
	router.SetTrustedProxies(nil)
	router.Run("localhost:8000");
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user != nil {
			c.Next()
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		}
	}
}

func getAccessToken(c *gin.Context) {
	godotenv.Load("../.env")
	CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("TWITCH_CLIENT_SECRET")
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", CLIENT_ID, CLIENT_SECRET)
	
	postBody, err := json.Marshal(&AccessToken{})
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

	var body AccessToken
	
	if err := json.Unmarshal(resBody, &body); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, body)
}

func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Parameters can't be empty"})
		return
	}

	// Placeholder for some sql logic that will be implemented later
	
	if username == "" && password == "" {
		session.Set("user", username)
		err := session.Save()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate session token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Successfully authenticated user"})
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	if user != nil {
		log.Println(user)
		session.Delete("user")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
	}
}

func status(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in"})
}

func private(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	c.JSON(http.StatusOK, gin.H{"status": user})
}
