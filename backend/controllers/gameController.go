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
	"time"

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

func GetIGDBRoute (c *gin.Context, url string, fields string, model *[]interface{}) {
	CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	bearer, err := c.Cookie("IGDBAccessToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
	}

	bearer = "Bearer " + bearer

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	reqBody := fmt.Sprintf("%s; where id = %s;", fields, string(jsonData))

	reader := bytes.NewReader([]byte(reqBody))

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Client-ID", CLIENT_ID)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	
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
			err := json.NewDecoder(res.Body).Decode(model)
			if err != nil {
				log.Fatal(err)
			}

			c.JSON(http.StatusOK, gin.H{
				"response": (*model)[0],
			})
	}
}

func GetGame(c *gin.Context) {
	// CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	// bearer, err := c.Cookie("IGDBAccessToken")
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"error": "Not authorized",
	// 	})
	// }

	// bearer = "Bearer " + bearer

	// jsonData, err := io.ReadAll(c.Request.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// url := "https://api.igdb.com/v4/games"
	// reqBody := fmt.Sprintf("fields id, aggregated_rating, aggregated_rating_count, alternative_names, collections, cover, game_modes, genres, involved_companies, name, platforms, first_release_date, slug, summary, url; where id = %s;", string(jsonData))

	// reader := bytes.NewReader([]byte(reqBody))

	// req, err := http.NewRequest("POST", url, reader)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// req.Header.Set("Client-ID", CLIENT_ID)
	// req.Header.Add("Authorization", bearer)

	// client := &http.Client{Timeout: 10 * time.Second}

	// res, err := client.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer res.Body.Close()

	// var game []models.Game
	
	// switch {
	// 	case res.StatusCode == 401:
	// 		c.JSON(http.StatusUnauthorized, gin.H{
	// 			"error": "Not authorized",
	// 		})
	// 	case res.StatusCode >= 400:
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "Bad request",
	// 		})
	// 	default:
	// 		err := json.NewDecoder(res.Body).Decode(&game)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}

	// 		c.JSON(http.StatusOK, gin.H{
	// 			"response": game[0],
	// 		})
	// }
	var game *[]models.Game
	GetIGDBRoute(
		c, 
		"https://api.igdb.com/v4/games", 
		"fields id, aggregated_rating, aggregated_rating_count, alternative_names, collections, cover, game_modes, genres, involved_companies, name, platforms, first_release_date, slug, summary, url", 
		game,
	)
}

func GetGameCover(c *gin.Context) {
	CLIENT_ID := os.Getenv("TWITCH_CLIENT_ID")
	bearer, err := c.Cookie("IGDBAccessToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authorized",
		})
	}

	bearer = "Bearer " + bearer

	jsonData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Fatal(err)
	}

	url := "https://api.igdb.com/v4/covers"
	reqBody := fmt.Sprintf("fields game, image_id; where id = %s;", string(jsonData))

	reader := bytes.NewReader([]byte(reqBody))

	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Client-ID", CLIENT_ID)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{Timeout: 10 * time.Second}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var game []models.Cover
	
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
