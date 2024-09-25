package main

import (
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
	skipQuery         = "SkipQuery"
)

var (
	searchButton       = "Search"
	searchMovieButton  = "Search Movie"
	searchTVButton     = "Search Series"
	searchPersonButton = "Search Person"
	cancelButton       = "Cancel"

	movieListsButton  = "Movie lists"
	movieListPopular  = "Popular"
	movieListTopRated = "Top rated"
	//movieListUpcoming = "Upcoming"
	tvListsButton     = "TV lists"
	personListsButton = "Person lists"
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

func createKeyboard(buttons ...string) tgbotapi.ReplyKeyboardMarkup {
	var keyboard []tgbotapi.KeyboardButton
	for _, button := range buttons {
		keyboard = append(keyboard, tgbotapi.NewKeyboardButton(button))
	}

	return tgbotapi.NewReplyKeyboard(keyboard)
}

func handleUserAction(update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	switch {
	//Start app
	case update.Message.Text == "/start":
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hi! Choose action")
		msg.ReplyMarkup = createKeyboard(searchButton, movieListsButton, tvListsButton, personListsButton)
		return msg, nil
	//Cancel action
	case update.Message.Text == cancelButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Choose action")
		msg.ReplyMarkup = createKeyboard(searchButton, movieListsButton, tvListsButton, personListsButton)
		return msg, nil
	//Search options
	case update.Message.Text == searchButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "What do you search?")
		msg.ReplyMarkup = createKeyboard(searchMovieButton, searchTVButton, searchPersonButton, cancelButton)
		return msg, nil
	//Movie lists options
	case update.Message.Text == movieListsButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a movie list")
		msg.ReplyMarkup = createKeyboard(movieListPopular, cancelButton)
		return msg, nil
	//Start of search movie flow
	case update.Message.Text == searchMovieButton:
		currentUserAction = searchMovieButton
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter movie title"), nil
	//Start of search TV flow
	case update.Message.Text == searchTVButton:
		currentUserAction = searchTVButton
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter series title"), nil
	//Start of search person flow
	case update.Message.Text == searchPersonButton:
		currentUserAction = searchPersonButton
		return tgbotapi.NewMessage(update.Message.Chat.ID, "Enter person's name"), nil
	//Search movie by title
	case currentUserAction == searchMovieButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		err := actionSearchMovie(update, &msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("search movie error: %v", err)
		}
		return msg, nil
	//Search TV by title
	case currentUserAction == searchTVButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

		err := actionSearchTV(update, &msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("search TV error: %v", err)
		}
		return msg, nil
	//Search person by name
	case currentUserAction == searchPersonButton:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		err := actionSearchPerson(update, &msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("search person error: %v", err)
		}
		return msg, nil
	//Popular movie list
	case update.Message.Text == movieListPopular:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyMarkup = createKeyboard(movieListPopular, cancelButton)
		err := actionPopularMovieList(&msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("popular movie list error: %v", err)
		}
		return msg, nil
	//Top rated movie list
	case update.Message.Text == movieListTopRated:
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyMarkup = createKeyboard(movieListTopRated, cancelButton)
		err := actionPopularMovieList(&msg)
		if err != nil {
			return tgbotapi.MessageConfig{}, fmt.Errorf("top rated movie list error: %v", err)
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
Rating: %.2f
Release date: %s
Overview: %s`,
		movie.Original_title, movie.Vote_average, movie.Release_date, movie.Overview)

	//fmt.Println(msg.Text)
	return nil
}

func actionSearchTV(update tgbotapi.Update, msg *tgbotapi.MessageConfig) error {
	path := "/search/tv"
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

	var tv TV
	if err = json.Unmarshal(body, &tv); err != nil {
		return fmt.Errorf("parse body response error: %v", err)
	}

	msg.Text = fmt.Sprintf(`Original title: %s
Rating: %.2f
Release date: %s
Overview: %s`,
		tv.OriginalName, tv.VoteAverage, tv.FirstAirDate, tv.Overview)

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

func actionPopularMovieList(msg *tgbotapi.MessageConfig) error {
	path := "/movie/popular"

	url, err := createURL(serverHost, path, skipQuery)
	if err != nil {
		return fmt.Errorf("create url error: %v", err)
	}

	// Make request to Server
	body, err := makeRequestToServer(url)
	if err != nil {
		return err
	}

	var popularMovies MoviesList
	if err = json.Unmarshal(body, &popularMovies); err != nil {
		return fmt.Errorf("parse body response error: %v", err)
	}
	msg.Text = `The most popular movies right now

`
	for i, movie := range popularMovies.Movies[:5] {
		msg.Text += fmt.Sprintf(`Top %d
Original title: %s
Rating: %.2f
Release date: %s
Overview: %s

`,
			i+1, movie.Original_title, movie.Vote_average, movie.Release_date, movie.Overview)
	}

	return nil
}

func actionTopRatedMovieList(msg *tgbotapi.MessageConfig) error {
	path := "/movie/top_rated"

	url, err := createURL(serverHost, path, skipQuery)
	if err != nil {
		return fmt.Errorf("create url error: %v", err)
	}

	// Make request to Server
	body, err := makeRequestToServer(url)
	if err != nil {
		return err
	}

	var popularMovies MoviesList
	if err = json.Unmarshal(body, &popularMovies); err != nil {
		return fmt.Errorf("parse body response error: %v", err)
	}
	msg.Text = `The most popular movies right now

`
	for i, movie := range popularMovies.Movies[:5] {
		msg.Text += fmt.Sprintf(`Top %d
Original title: %s
Rating: %.2f
Release date: %s
Overview: %s

`,
			i+1, movie.Original_title, movie.Vote_average, movie.Release_date, movie.Overview)
	}

	return nil
}

func createURL(host, path, query string) (url.URL, error) {
	// Create full url with query param
	params := url.Values{}

	if query != skipQuery {
		params.Add("query", query)
	}

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
	fmt.Println("     Body make request:", string(body))
	if err != nil {
		return nil, fmt.Errorf("read server response error: %v", err)
	}
	return body, nil
}
