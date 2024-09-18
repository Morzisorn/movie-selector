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
		ID                int64   `json:"id"`
		Original_title    string  `json:"original_title"`
		Genre_ids         []int64 `json:"genre_ids"`
		Original_language string  `json:"original_language"`
		Overview          string  `json:"overview"`
		Release_date      string  `json:"release_date"`
		Vote_average      float64 `json:"vote_average"`
	}
)

func getTMDBMovies(c *fiber.Ctx) error {
	/*
		var resp GetTgMoviesRequest
		if err := c.BodyParser(&resp); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	*/
	fmt.Println("Enter getTMDBMovies")
	query := c.Params("query")
	if query == "" {
		fmt.Println("No query in params")
		fmt.Println(c.Request())
		fmt.Println(c.AllParams())
		return c.Status(fiber.StatusBadRequest).SendString("No query in params")
	}

	tmdbMoviesResponse, err := SearchMovie(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Use SearchMovie error: %w", err))
	}
	var tmdbMovies GetTMBDMoviesResponse

	err = json.Unmarshal([]byte(tmdbMoviesResponse), &tmdbMovies)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	fmt.Println("Finish getTMDBMovies")
	return c.JSON(tmdbMovies.Movies[0].Release_date)
}

func SearchMovie(query string) (string, error) {
	fmt.Println("Enter Search Movie")
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

	fmt.Println("Finish Search Movie")
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
