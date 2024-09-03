package controllers

import (
	"backend/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetAccessToken(c *gin.Context) {
	CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("TWITCH_CLIENT_SECRET")
	url := fmt.Sprintf("https://id.twitch.tv/oauth2/token?client_id=%s&client_secret=%s&grant_type=client_credentials", CLIENT_ID, CLIENT_SECRET)

	postBody, err := json.Marshal(&models.AccessToken{})
	if err != nil {
		log.Fatal(err)
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

	c.SetCookie("IGDBAccessToken", body.AccessToken, 2592000, "", "", false, true)
	c.JSON(http.StatusOK, body)
}

func GetGame(c *gin.Context) {
	token, err := c.Cookie("IGDBAccessToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
	}

	url := "https://api.igdb.com/v4/games"
	jsonData, err := io.ReadAll(c.Request.Body)
	body := fmt.Sprintf("fields id, aggregated_rating, aggregated_rating_count, alternative_names, collections, cover, game_modes, genres, involved_companies, name, platforms, first_release_date, slug, summary, url; where id = %s;", string(jsonData))

	postBody, err := json.Marshal(&models.Game{})
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(postBody)

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Client-ID", token)

}
