package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/schafdog/gameapi/cassandra"
	"net/http"
)

func ParsePostUserRequest(r *http.Request) (user User, err error) {
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&user)
	if err != nil {
		fmt.Println("Parse Post User Request errors", err.Error())
	}
	return user, err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var err error
	user, err = ParsePostUserRequest(r)
	if err != nil {
		HandleHttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = persist(&user)
	if err != nil {
		HandleHttpResponse(w, http.StatusInternalServerError, "Failed to persist "+user.Id.String()+": "+err.Error())
		return
	}
	HandleNewUserResponse(w, user)
}

func persist(user *User) error {
	var gocqlUuid gocql.UUID

	// generate a unique UUID for this user
	gocqlUuid = gocql.TimeUUID()
	fmt.Println("creating a new user", gocqlUuid, " for ", user.Name)
	user.Id = gocqlUuid
	// write data to Cassandra
	err := cassandra.Session.Query(`INSERT INTO user (id, name, score, gamesplayed) VALUES (?, ?, ?, ?)`, gocqlUuid, user.Name, 0, 0).Exec()
	return err
}

func HandleNewUserResponse(w http.ResponseWriter, user User) {
	fmt.Println("user_id", user.Id)
	decoder := json.NewEncoder(w)
	decoder.SetIndent("", "   ")
	decoder.Encode(NewUserResponse{Id: user.Id, Name: user.Name})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []NewUserResponse
	m := map[string]interface{}{}

	query := "SELECT id,name FROM User"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		fmt.Printf("User{ id: %v, name: %v }", m["id"], m["name"])
		userList = append(userList, NewUserResponse{
			Id:   m["id"].(gocql.UUID),
			Name: m["name"].(string),
		})
		m = map[string]interface{}{}
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "   ")
	encoder.Encode(UsersResponse{Users: userList})
}

func HandleErrorsResponse(w http.ResponseWriter, errs []string) {
	fmt.Println("errors", errs)
	json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
}
