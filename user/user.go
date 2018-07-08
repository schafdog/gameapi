package User

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/schafdog/gameapi/cassandra"
	"net/http"
)

type userCreateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NewUserRequest struct {
	Name *string `json:"name"`
}

// UserResponse builds a payload of new user resource ID
type UserResponse struct {
	Id   gocql.UUID `json:"id"`
	Name string     `json:"name"`
}

// UsersResponse to form payload of an array of User structs
type UsersResponse struct {
	Users []UserResponse `json:"users"`
}

func ParsePostUserRequest(r *http.Request) (user User, err error) {
	var newUser NewUserRequest
	var uuidStr string
	var uuid gocql.UUID
	uuid = gocql.TimeUUID()
	uuidStr = r.Header.Get("X-UUID")
	if uuidStr != "" {
		fmt.Printf("Suggestion for UUID: %v\n", uuidStr)
		uuid, err = gocql.ParseUUID(uuidStr)
		if err != nil {
			fmt.Printf("Failed to parse X-UUID: %v\n", uuidStr)
		}
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUser)
	if err != nil {
		fmt.Println("Parse Post User Request errors", err.Error())
		return User{}, err
	}
	if newUser.Name == nil || len(*newUser.Name) == 0 {
		return User{}, errors.New("User: Name is missing or empty")
	}
	return User{Name: *newUser.Name, Id: &uuid}, err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user User
	var err error
	user, err = ParsePostUserRequest(r)
	if err != nil {
		HandleHttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	err = Persist(&user)
	if err != nil {
		HandleHttpResponse(w, http.StatusInternalServerError, "Failed to persist "+user.Id.String()+": "+err.Error())
		return
	}
	HandleNewUserResponse(w, user)
}

func Persist(user *User) error {
	var gocqlUuid gocql.UUID

	// generate a unique UUID for this user
	if user.Id == nil {
		gocqlUuid = gocql.TimeUUID()
		fmt.Printf("creating a new userid %v for %v\n", gocqlUuid, user.Name)
		user.Id = &gocqlUuid
	} else {
		fmt.Printf("Using suggestion %v from header for %v\n", user.Id, user.Name)
	}
	// write data to Cassandra
	var friends []gocql.UUID
	err := cassandra.Session.Query(`INSERT INTO user (id, name, score, gamesplayed, friends) VALUES (?, ?, ?, ?, ?)`, user.Id, user.Name, 0, 0, friends).Exec()
	return err
}

func Delete(uuid gocql.UUID) error {
	// write data to Cassandra
	err := cassandra.Session.Query(`DELETE FROM user where id = ?`, uuid).Exec()
	return err
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var uuid, error = ParseUserid(r)
	if error != nil {
		fmt.Println("Failed to parse user id")
		HandleHttpResponse(w, http.StatusBadRequest, error.Error())
		return
	}
	fmt.Printf("userid %v", uuid)
	error = Delete(uuid)
	if error != nil {
		fmt.Printf("Failed to delete user %v: %v\n", uuid, error.Error())
		// Handle not found and internal server error
		HandleHttpResponse(w, http.StatusInternalServerError, error.Error())
	}

}
func HandleNewUserResponse(w http.ResponseWriter, user User) {
	fmt.Println("user_id", user.Id)
	decoder := json.NewEncoder(w)
	decoder.SetIndent("", "   ")
	decoder.Encode(UserResponse{Id: *user.Id, Name: user.Name})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var userList []UserResponse
	m := map[string]interface{}{}

	query := "SELECT id,name FROM User"
	iterable := cassandra.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		fmt.Printf("User{ id: %v, name: %v }\n", m["id"], m["name"])
		userList = append(userList, UserResponse{
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
