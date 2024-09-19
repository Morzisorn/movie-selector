package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var (
	Tmdb_api_key string
	Tg_api_key   string
)

func getEnvKeys() {
	err := godotenv.Load()
	if err != nil {
		logrus.Fatal("Ошибка загрузки файла .env")
	}

	// Получаем API ключ из переменных окружения
	Tmdb_api_key = os.Getenv("Tmdb_api_key")
	if Tmdb_api_key == "" {
		fmt.Println("API_KEY не установлен")
	}

	Tg_api_key = os.Getenv("Tg_api_key")
	if Tmdb_api_key == "" {
		fmt.Println("API_KEY не установлен")
	}
}

func main() {
	getEnvKeys()
	go StartServer()
	StartBot()
	fmt.Println("Server and bot started")
}
