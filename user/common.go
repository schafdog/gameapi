package User

// UserState to form payload returning a single User struct
type UserState struct {
	GamesPlayed int `json:"gamesPlayed"`
	Highscore   int `json:"score"`
}

// ErrorResponse returns an array of error strings if appropriate
type ErrorResponse struct {
	Error error `json:"error"`
}
