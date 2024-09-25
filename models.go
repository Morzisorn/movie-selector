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

	MoviesList struct {
		Movies []Movie
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

	TV struct {
		Adult            bool     `json:"adult,omitempty"`
		BackdropPath     string   `json:"backdrop_path,omitempty"`
		GenreIDs         []int    `json:"genre_ids,omitempty"`
		ID               int      `json:"id,omitempty"`
		OriginCountry    []string `json:"origin_country,omitempty"`
		OriginalLanguage string   `json:"original_language,omitempty"`
		OriginalName     string   `json:"original_name"`
		Overview         string   `json:"overview"`
		Popularity       float64  `json:"popularity,omitempty"`
		PosterPath       string   `json:"poster_path,omitempty"`
		FirstAirDate     string   `json:"first_air_date,omitempty"`
		Name             string   `json:"name,omitempty"`
		VoteAverage      float64  `json:"vote_average,omitempty"`
		VoteCount        int      `json:"vote_count,omitempty"`
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
