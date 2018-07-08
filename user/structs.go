package User

import (
	"github.com/gocql/gocql"
)

// User struct to hold profile data for our user
type User struct {
	Id          gocql.UUID   `json:"id"`
	Name        string       `json:"name"`
	Highscore   int          `json:"score"`
	GamesPlayed int          `json:"gamesPlayed"`
	Friends     []gocql.UUID `json:"friends"`
}

// UserStat to form payload returning a single User struct
type Stat struct {
	Id          gocql.UUID `json:"id"`
	Highscore   int        `json:"score"`
	GamesPlayed int        `json:"gamesPlayed"`
}

// UserStat to form payload returning a single User struct
type UserStat struct {
	GamesPlayed int `json:"gamesPlayed"`
	Highscore   int `json:"score"`
}

// UsersResponse to form payload of an array of User structs
type UsersResponse struct {
	Users []NewUserResponse `json:"users"`
}

// NewUserResponse builds a payload of new user resource ID
type NewUserResponse struct {
	Id   gocql.UUID `json:"id"`
	Name string     `json:"name"`
}

// ErrorResponse returns an array of error strings if appropriate
type ErrorResponse struct {
	Errors []string `json:"errors"`
}
