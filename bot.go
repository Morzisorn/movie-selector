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

var (
	serverHost        = "http://localhost:3000"
	currentUserAction string
)

var (
	searchMovie  = "Search Movie"
	searchPerson = "Search Person"
)

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

func handleTgUpdates(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	if update.Message != nil { // Check new messages
		msg, err := handleUserAction(update)
		if err != nil {
			msg.Text = fmt.Sprint(err.Error())
		}

		bot.Send(msg)
	}
}

func handleUserAction(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	switch {
	case update.Message.Text == "/start": //Start app
		buttonSearchMovie := tgbotapi.NewKeyboardButton(searchMovie)
		buttonSearchActor := tgbotapi.NewKeyboardButton(searchPerson)
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(buttonSearchMovie, buttonSearchActor),
		)

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hi! Choose action")
		msg.ReplyMarkup = keyboard
		return msg, nil
	case update.Message.Text == searchMovie: //Start of search movie flow
		currentUserAction = searchMovie
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter movie title"), nil
	case update.Message.Text == searchPerson: //Start of search person flow
		currentUserAction = searchPerson
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter person's name"), nil
	case currentUserAction == searchMovie: //Search movie by title
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		err := actionSearchMovie(update, &msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("search movie error: %v", err)
		}
		return msg, nil
	case currentUserAction == searchPerson: //Search person by name
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		err := actionSearchPerson(update, &msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("search person error: %v", err)
		}
		return msg, nil
	default:
		return tgbotapi.NewMessage(update.Message.Chat.ID, "I'm broken"), nil
	}
}

func actionSearchMovie(update tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	path := "/search/movie"
	query := update.Message.Text

	url, err := createURL(serverHost, path, query)
	if err != nil {
		return fmt.Errorf("create url error: %v", err)
	}

	// Make request to Server
	body, err := makeRequestToServer(url)
	if err != nil {
		return err
	}

	var movie Movie
	if err = json.Unmarshal(body, &movie); err != nil {
		return fmt.Errorf("parse body response error: %v", err)
	}

	msg.Text = fmt.Sprintf(`Original title: %s
Release date: %s
Overview: %s`,
		movie.Original_title, movie.Release_date, movie.Overview)

	//fmt.Println(msg.Text)
	return nil
}

func actionSearchPerson(update tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	path := "/search/person"
	query := update.Message.Text

	url, err := createURL(serverHost, path, query)
	if err != nil {
		return fmt.Errorf("create url error: %v", err)
	}

	// Make request to Server
	body, err := makeRequestToServer(url)
	if err != nil {
		return err
	}

	var person Person
	if err = json.Unmarshal(body, &person); err != nil {
		return fmt.Errorf("parse body response error: %v", err)
	}
	var personKnownFor = person.KnownFor[0].OriginalTitle

	for i := 1; i < len(person.KnownFor) && i <= 3; i++ {
		personKnownFor += ", " + person.KnownFor[i].OriginalTitle
	}

	msg.Text = fmt.Sprintf(`%s
%s
Known for: %s`,
		person.Name, person.KnownForDepartment, personKnownFor)

	return nil
}

func createURL(host, path, query string) (url.URL, error) {
	// Create full url with query param
	params := url.Values{}
	params.Add("query", query)

	baseURL, err := url.Parse(host)
	if err != nil {
		return url.URL{}, fmt.Errorf("host parsing error: %v", err)
	}

	baseURL.Path += path
	baseURL.RawQuery = params.Encode()
	return *baseURL, nil
}

func makeRequestToServer(url url.URL) ([]byte, error) {
	response, err := http.Get(url.String())
	if err != nil {
		return nil, fmt.Errorf("request server error: %v", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("read server response error: %v", err)
	}
	return body, nil
}
