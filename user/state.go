package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/schafdog/gameapi/db"
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

func ParseStateRequest(r *http.Request) (db.State, error) {
	var state db.State
	uuid, error := ParseUserid(r)
	if error != nil {
		return state, error
	} else {
		state.Id = uuid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&state)
		if err != nil {
			fmt.Println("Failed to parse ParseStateRequest", err)
			return state, err
		}
	}
	fmt.Printf("STAT: %v\n", state)
	return state, error
}

func GetState(w http.ResponseWriter, r *http.Request) {
	var state *db.State
	var uuid, error = ParseUserid(r)
	if error != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: error})
		return
	}
	fmt.Printf("Found uuid %v in url \n", uuid)
	state, error = DB.GetState(uuid)
	if error != nil {
		msg := fmt.Sprintf("Failed to lookup uuid %v in DB \n", uuid)
		fmt.Printf(msg)
		HandleHttpResponse(w, http.StatusNotFound, msg)
		return
	}
	handleStateResponse(w, *state, error)
}

func PutState(w http.ResponseWriter, r *http.Request) {
	var state db.State
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	uuid, error = gocql.ParseUUID(vars["userid"])
	fmt.Printf("User Id: %s\n", uuid)
	if error != nil {
		handleStateResponse(w, state, error)
		return
	}
	state, error = ParseStateRequest(r)
	if error != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: error})
		return
	}
	state.Id = uuid
	fmt.Printf("Userid %v Played %v games with highscore %v \n", state.Id, state.GamesPlayed, state.Highscore)
	error = DB.SetState(uuid, state)
	handlePutResponse(w, state, error)
}

func handlePutResponse(w http.ResponseWriter, state db.State, err error) {
	if err == nil {
		fmt.Println("user_id", state.Id)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("error: ", err)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err})
	}
}

func handleStateResponse(w http.ResponseWriter, state db.State, err error) {
	if err == nil {
		fmt.Println("user id", state.Id, state.GamesPlayed, state.Highscore)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "   ")
		encoder.Encode(UserState{GamesPlayed: state.GamesPlayed, Highscore: state.Highscore})
	} else {
		fmt.Println("error ", err)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err})
	}
}
