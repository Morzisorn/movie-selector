package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
)

var (
	baseURL = "https://api.themoviedb.org/3"
)

type (
	GetTgMoviesRequest struct {
		MoviesRequest string
	}

	GetTMBDMoviesResponse struct {
		Movies []Movie `json:"results"`
	}

	Movie struct {
		ID                int64   `json:"id,omitempty"`
		Original_title    string  `json:"original_title,omitempty"`
		Genre_ids         []int64 `json:"genre_ids,omitempty"`
		Original_language string  `json:"original_language,omitempty"`
		Overview          string  `json:"overview,omitempty"`
		Release_date      string  `json:"release_date,omitempty"`
		Vote_average      float64 `json:"vote_average,omitempty"`
	}
)

func getTMDBMovies(c *fiber.Ctx) error {
	query := c.Query("query")
	if query == "" {
		fmt.Println("No query in params")
		return c.Status(fiber.StatusBadRequest).SendString("No query in params")
	}

	tmdbMoviesResponse, err := SearchMovie(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Use SearchMovie error: %v", err))
	}
	var tmdbMovies GetTMBDMoviesResponse

	err = json.Unmarshal([]byte(tmdbMoviesResponse), &tmdbMovies)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp := Movie{
		Original_title: tmdbMovies.Movies[0].Original_title,
		Release_date:   tmdbMovies.Movies[0].Release_date,
		Overview:       tmdbMovies.Movies[0].Overview,
	}

	return c.JSON(resp)
}

func SearchMovie(query string) (string, error) {
	url := baseURL + "/search/movie" + "?query=" + query

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("create get request to tmdb: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+Tmdb_api_key)
	req.Header.Add("accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do get to tmdb: %w", err)
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return string(body), nil

}

func StartServer() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Get("/search/movie", getTMDBMovies)

	logrus.Fatal(app.Listen(":3000"))
}
