package main

type (
	Movie struct {
		ID                int64   `json:"id,omitempty"`
		Original_title    string  `json:"original_title"`
		Genre_ids         []int64 `json:"genre_ids,omitempty"`
		Original_language string  `json:"original_language,omitempty"`
		Overview          string  `json:"overview"`
		Release_date      string  `json:"release_date"`
		Vote_average      float64 `json:"vote_average,omitempty"`
	}

	Person struct {
		Adult              bool           `json:"adult,omitempty"`
		Gender             int            `json:"gender,omitempty"`
		ID                 int            `json:"id,omitempty"`
		KnownForDepartment string         `json:"known_for_department"`
		Name               string         `json:"name"`
		OriginalName       string         `json:"original_name"`
		Popularity         float64        `json:"popularity,omitempty"`
		ProfilePath        string         `json:"profile_path,omitempty"`
		KnownFor           []KnownForItem `json:"known_for"`
	}

	KnownForItem struct {
		BackdropPath     string  `json:"backdrop_path,omitempty"`
		ID               int     `json:"id,omitempty"`
		Title            string  `json:"title"`
		OriginalTitle    string  `json:"original_title"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path,omitempty"`
		MediaType        string  `json:"media_type"`
		Adult            bool    `json:"adult,omitempty"`
		OriginalLanguage string  `json:"original_language,omitempty"`
		GenreIDs         []int   `json:"genre_ids"`
		Popularity       float64 `json:"popularity,omitempty"`
		ReleaseDate      string  `json:"release_date"`
		Video            bool    `json:"video,omitempty"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count,omitempty"`
	}
)
