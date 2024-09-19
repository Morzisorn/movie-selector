package main

import (
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

func handleTgUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil { // Check new messages
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		// Create full url with query params
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

		// Make request to Server
		response, err := http.Get(baseURL.String())
		if err != nil {
			logrus.Println("Error making request to server:", err)
			msg.Text = "Error contacting the server"
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)

		if err != nil {
			fmt.Println("ReadAll error:", err)
			msg.Text = "Error contacting the server"
		}

		var movie Movie
		if err = json.Unmarshal(body, &movie); err != nil {
			logrus.Println("Error making request to server:", err)
			fmt.Printf("Error contacting the server: Unmarshal %s", err)
		}

		msg.Text = fmt.Sprintf(`Original title: %s
Release date: %s
Overview: %s`,
			movie.Original_title, movie.Release_date, movie.Overview)

		//fmt.Println(msg.Text)
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
