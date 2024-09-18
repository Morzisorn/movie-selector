package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func handleTgUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil { // Проверяем, есть ли сообщение
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// Здесь вы можете обработать текст сообщения и отправить запрос к вашему серверу
		host := "http://localhost:3000"
		path := "/search/movie"
		params := url.Values{}
		params.Add("query", update.Message.Text)

		baseURL, err := url.Parse(host)
		if err != nil {
			fmt.Println("Host parsing error:", err)
		}

		baseURL.Path += path
		baseURL.RawQuery = params.Encode()

		fmt.Printf("baseURL: %s\n", baseURL.String())
		response, err := http.Get(baseURL.String()) // Пример запроса к вашему серверу
		if err != nil {
			logrus.Println("Error making request to server:", err)
			msg.Text = "Error contacting the server"
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logrus.Println("Error making request to server:", err)
			msg.Text = "Error contacting the server"
		}
		msg.Text = string(body)
		bot.Send(msg)
	}
}

func StartBot() {
	bot, err := tgbotapi.NewBotAPI(Tg_api_key)
	if err != nil {
		logrus.Panic(err)
	}

	bot.Debug = true
	logrus.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		handleTgUpdates(bot, update) //was go handle...
	}
}
