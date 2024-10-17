package controllers

import (
	"backend/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetGame(c *gin.Context) {
	API_KEY := os.Getenv("API_KEY")
	queryString := c.Query("id")

	url := fmt.Sprintf("https://www.giantbomb.com/api/game/%s/?api_key=%s&format=json", queryString, API_KEY)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var game []models.Game
	
	switch {
		case res.StatusCode == 401:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
		case res.StatusCode >= 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})
		default:
			err := json.NewDecoder(res.Body).Decode(&game)
			if err != nil {
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"response": game[0],
			})
	}
}

func SearchGames(c *gin.Context) {
	API_KEY := os.Getenv("API_KEY")
	search := c.Query("search")

	url := fmt.Sprintf("https://www.giantbomb.com/api/search/?api_key=%s&format=json&query=%s", API_KEY, search)

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var game []models.Game
	
	switch {
		case res.StatusCode == 401:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Not authorized",
			})
		case res.StatusCode >= 400:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Bad request",
			})
		default:
			err := json.NewDecoder(res.Body).Decode(&game)
			if err != nil {
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"response": game,
			})
	}
}

func WriteGame(c *gin.Context) {

}