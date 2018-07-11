package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"net/http"
)

// UserState to form payload returning a single User struct
type UserState struct {
	GamesPlayed int `json:"gamesPlayed"`
	Highscore   int `json:"score"`
}

// ErrorResponse returns an array of error strings if appropriate
type ErrorResponse struct {
	Error error `json:"error"`
}

func ParseUserid(r *http.Request) (gocql.UUID, error) {
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	fmt.Printf("User Id: %v\n", vars["userid"])
	uuid, error = gocql.ParseUUID(vars["userid"])
	return uuid, error
}

func HandleHttpResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if message != "" || len(message) > 0 {
		json.NewEncoder(w).Encode(message)
	}
}
