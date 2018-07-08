package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/schafdog/gameapi/cassandra"
	"net/http"
)

func ParseUserRequest(r *http.Request) (User, []string) {
	var user User
	var errs []string
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil {
		fmt.Println("errors", errs)
		//		w.WriteHeader(http.StatusBadRequest)
		//		json.NewEncoder(w).Encode(ErrorResponse{Errors: err})
	}
	return user, errs
}

func ParseUserid(r *http.Request) (gocql.UUID, error) {
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	fmt.Printf("User Id: %v\n", vars["userid"])
	uuid, error = gocql.ParseUUID(vars["userid"])
	return uuid, error
}

func ParseStatRequest(r *http.Request) (Stat, []string) {
	var stat Stat
	var errs []string
	uuid, error := ParseUserid(r)
	if error != nil {
		errs = append(errs, error.Error())
	} else {
		stat.Id = uuid
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&stat)
		if err != nil {
			fmt.Println("Failed to parse ParseStatRequest", err)
			errs = append(errs, err.Error())
		}
	}
	fmt.Printf("STAT: %v\n", stat)
	return stat, errs
}

func GetStat(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var stat Stat
	var uuid, error = ParseUserid(r)
	if error != nil {
		errs = append(errs, "Failed to parse uuid in: "+r.URL.Path)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
		return
	}
	stat.Id = uuid
	errs = getStat(&stat)
	handleStatResponse(w, stat, errs)
}

func lookup(stat *Stat) []string {
	var result []string
	return result
}

func PutStat(w http.ResponseWriter, r *http.Request) {
	var stat Stat
	var errs []string
	var uuid gocql.UUID
	var error error
	vars := mux.Vars(r)
	uuid, error = gocql.ParseUUID(vars["userid"])
	fmt.Printf("User Id: %s\n", uuid)
	if error != nil {
		handleStatResponse(w, stat, errs)
		return
	}
	stat, errs = ParseStatRequest(r)
	if len(errs) != 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
		return
	}
	stat.Id = uuid
	fmt.Printf("Userid %v Played %v games with highscore %v \n", stat.Id, stat.GamesPlayed, stat.Highscore)
	errs = persistStat(&stat)
	handlePutResponse(w, stat, errs)
}

func persistStat(stat *Stat) []string {
	var errs []string
	if err := cassandra.Session.Query(`
      UPDATE user set gamesPlayed = ?, score = ? where id = ?`,
		stat.GamesPlayed, stat.Highscore, stat.Id).Exec(); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}

func getStat(stat *Stat) []string {
	var errs []string
	if err := cassandra.Session.Query(`
      select gamesPlayed, score from User where id = ?`,
		stat.Id).Scan(&stat.GamesPlayed, &stat.Highscore); err != nil {
		errs = append(errs, err.Error())
	}
	return errs
}

func handlePutResponse(w http.ResponseWriter, stat Stat, errs []string) {
	if errs == nil || len(errs) == 0 {
		fmt.Println("user_id", stat.Id)
		w.WriteHeader(http.StatusOK)
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

func handleStatResponse(w http.ResponseWriter, stat Stat, errs []string) {
	if errs == nil {
		fmt.Println("user id", stat.Id, stat.GamesPlayed, stat.Highscore)
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "   ")
		encoder.Encode(UserStat{GamesPlayed: stat.GamesPlayed, Highscore: stat.Highscore})
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}
