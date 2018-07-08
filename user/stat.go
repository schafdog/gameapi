package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/schafdog/gameapi/cassandra"
	"net/http"
)

func ParseUserid(r *http.Request) (gocql.UUID, error) {
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	fmt.Printf("User Id: %v\n", vars["userid"])
	uuid, error = gocql.ParseUUID(vars["userid"])
	return uuid, error
}

func ParseStateRequest(r *http.Request) (State, []string) {
	var stat State
	var errs []string
	uuid, error := ParseUserid(r)
	if error != nil {
		errs = append(errs, error.Error())
	} else {
		stat.Id = uuid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&stat)
		if err != nil {
			fmt.Println("Failed to parse ParseStateRequest", err)
			errs = append(errs, err.Error())
		}
	}
	fmt.Printf("STAT: %v\n", stat)
	return stat, errs
}

func GetState(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var stat State
	var uuid, error = ParseUserid(r)
	if error != nil {
		errs = append(errs, "Failed to parse uuid in: "+r.URL.Path)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
		return
	}
	stat.Id = uuid
	errs = getState(&stat)
	handleStateResponse(w, stat, errs)
}

func lookup(stat *State) []string {
	var result []string
	return result
}

func PutState(w http.ResponseWriter, r *http.Request) {
	var stat State
	var errs []string
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	uuid, error = gocql.ParseUUID(vars["userid"])
	fmt.Printf("User Id: %s\n", uuid)
	if error != nil {
		handleStateResponse(w, stat, errs)
		return
	}
	stat, errs = ParseStateRequest(r)
	if len(errs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
		return
	}
	stat.Id = uuid
	fmt.Printf("Userid %v Played %v games with highscore %v \n", stat.Id, stat.GamesPlayed, stat.Highscore)
	errs = persistState(&stat)
	handlePutResponse(w, stat, errs)
}

func persistState(stat *State) []string {
	var errs []string
	if err := cassandra.Session.Query(`
      UPDATE user set gamesPlayed = ?, score = ? where id = ?`,
		stat.GamesPlayed, stat.Highscore, stat.Id).Exec(); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}

func getState(state *State) []string {
	var errs []string
	if err := cassandra.Session.Query(`
      select gamesPlayed, score from User where id = ?`,
		state.Id).Scan(&state.GamesPlayed, &state.Highscore); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}

func handlePutResponse(w http.ResponseWriter, state State, errs []string) {
	if errs == nil || len(errs) == 0 {
		fmt.Println("user_id", state.Id)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

func handleStateResponse(w http.ResponseWriter, state State, errs []string) {
	if errs == nil {
		fmt.Println("user id", state.Id, state.GamesPlayed, state.Highscore)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "   ")
		encoder.Encode(UserState{GamesPlayed: state.GamesPlayed, Highscore: state.Highscore})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
