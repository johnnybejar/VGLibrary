package controllers

import (
	"backend/constants"
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

func GetGame(c *gin.Context) {
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

	url := "https://api.igdb.com/v4/games"
	reqBody := fmt.Sprintf("fields id, aggregated_rating, aggregated_rating_count, collections, cover, game_modes, genres, involved_companies, name, platforms, first_release_date, slug, summary, url; where id = %s;", string(jsonData))

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

func SearchGames(c *gin.Context) {
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

	url := "https://api.igdb.com/v4/games"
	reqBody := fmt.Sprintf("fields id, aggregated_rating, aggregated_rating_count, collections, cover, game_modes, genres, involved_companies, name, platforms, first_release_date, slug, summary, url; search \"%s\";", string(jsonData))

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
	var game models.Game

	if c.Bind(&game) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to read body",
		})

		return
	}
	
	// Here, we have to make a few API calls to igdb to get more info about the game
	// This is to avoid having to make potentially hundreds of calls later
	// TODO: Some games will have empty fields, so handle those cases

	// 1. Collection
	var collection models.Collections = WriteGameHelper(c, "collections", "name").(models.Collections)
	fmt.Println(collection)

	// 2. Cover
	var cover models.Cover = WriteGameHelper(c, "covers", "game, image_id").(models.Cover)
	fmt.Println(cover)

	// 3. Game Modes
	gameModes := constants.GetGameModes()
	var gameModesString []string

	for i := 0; i < len(game.GameModes); i++ {
		gameModesString = append(gameModesString, gameModes[i])
	}
	fmt.Println(gameModesString)

	// 4. Genres
	genres := constants.GetGenres()
	var genresString []string
	
	for i := 0; i < len(game.Genres); i++ {
		genresString = append(genresString, genres[i])
	}
	fmt.Println(genresString)

	// 5. Companies (2 calls)


	// 6. Platforms

	fmt.Println(game)

	// res := initializers.DB.Create(&game)
	// if res.Error != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"Error": "Failure to create game",
	// 		"Code": res.Error,
	// 	})

	// 	initializers.DB.Delete("")
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"game": game,
	})
}

func WriteGameHelper(c *gin.Context, route string, fields string) interface{} {
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

	reqBody := fmt.Sprintf("fields %s; where id = %s;", fields, string(jsonData))

	reader := bytes.NewReader([]byte(reqBody))

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.igdb.com/v4/%s", route), reader)
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

	var obj []interface{}
	
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
			err := json.NewDecoder(res.Body).Decode(&obj)
			if err != nil {
				log.Fatal(err)
			}
	}

	return obj[0]
}
