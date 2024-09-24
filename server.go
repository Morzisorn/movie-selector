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
	tmdbBaseURL = "https://api.themoviedb.org/3"
)

type (
	TMDBClient struct{}

	MovieSearcher interface {
		SearchMovie(query string) (string, error)
	}

	PersonSearcher interface {
		SearchPerson(query string) (string, error)
	}

	GetTgMoviesRequest struct {
		MoviesRequest string
	}

	GetTMBDMoviesResponse struct {
		Movies []Movie `json:"results"`
	}

	GetTgPersonsRequest struct {
		PersonRequest string
	}

	GetTMBDPersonResponse struct {
		Persons []Person `json:"results"`
	}
)

func StartServer() {
	app := fiber.New(fiber.Config{
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Fiber!")
	})

	app.Get("/search/movie", func(c *fiber.Ctx) error {
		return getTMDBMovies(c, TMDBClient{})
	})
	app.Get("/search/person", func(c *fiber.Ctx) error {
		return getTMDBPerson(c, TMDBClient{})
	})

	logrus.Fatal(app.Listen(":3000"))
}

func getTMDBMovies(c *fiber.Ctx, s MovieSearcher) error {
	query := c.Query("query")
	if query == "" {
		fmt.Println("No query in params")
		return c.Status(fiber.StatusBadRequest).SendString("No query in params")
	}

	tmdbMoviesResponse, err := s.SearchMovie(query)
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

func getTMDBPerson(c *fiber.Ctx, s PersonSearcher) error {
	query := c.Query("query")
	if query == "" {
		fmt.Println("No query in params")
		return c.Status(fiber.StatusBadRequest).SendString("No query in params")
	}

	tmdbPersonResponse, err := s.SearchPerson(query)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("Use SearchPerson error: %v", err))
	}
	var tmdbPersons GetTMBDPersonResponse

	err = json.Unmarshal([]byte(tmdbPersonResponse), &tmdbPersons)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	resp := Person{
		Name:               tmdbPersons.Persons[0].Name,
		KnownForDepartment: tmdbPersons.Persons[0].KnownForDepartment,
		KnownFor:           tmdbPersons.Persons[0].KnownFor,
	}

	return c.JSON(resp)
}

func (t TMDBClient) SearchMovie(query string) (string, error) {
	path := "/search/movie"

	body, err := makeTMDBRequest(path, query)
	if err != nil {
		return "", err
	}

	return body, nil
}

func (t TMDBClient) SearchPerson(query string) (string, error) {
	path := "/search/person"

	body, err := makeTMDBRequest(path, query)
	if err != nil {
		return "", err
	}

	return body, nil
}

func makeTMDBRequest(path, query string) (string, error) {
	url, err := createURL(tmdbBaseURL, path, query)
	if err != nil {
		return "", fmt.Errorf("create url error: %v", err)
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return "", fmt.Errorf("create get request to tmdb: %v", err)
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
