package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	Tmdb_api_key string
	Tg_api_key   string
)

func getEnvKeys() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Load .env error")
	}

	// Получаем API ключ из переменных окружения
	Tmdb_api_key = os.Getenv("Tmdb_api_key")
	if Tmdb_api_key == "" {
		fmt.Println("API_KEY isn't set")
	}

	Tg_api_key = os.Getenv("Tg_api_key")
	if Tmdb_api_key == "" {
		fmt.Println("API_KEY isn't set")
	}
}

func main() {
	getEnvKeys()
	go StartServer()
	StartBot()
	fmt.Println("Server and bot started")
}
