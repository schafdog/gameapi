package User

import (
	"github.com/gocql/gocql"
)

// Model to hold profile data for our user
type User struct {
	Id          *gocql.UUID  `json:"id"`
	Name        string       `json:"name"`
	Highscore   int          `json:"score"`
	GamesPlayed int          `json:"gamesPlayed"`
	Friends     []gocql.UUID `json:"friends"`
}

// State to form payload returning a single User state
type State struct {
	Id          gocql.UUID `json:"id"`
	Highscore   int        `json:"score"`
	GamesPlayed int        `json:"gamesPlayed"`
}

// UserStat to form payload returning a single User struct
type UserState struct {
	GamesPlayed int `json:"gamesPlayed"`
	Highscore   int `json:"score"`
}

// ErrorResponse returns an array of error strings if appropriate
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
