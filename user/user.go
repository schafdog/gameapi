package User

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/schafdog/gameapi/db"
	"net/http"
)

var DB db.UserDatabase

func InitDB(newDB db.UserDatabase) {
	DB = newDB
}

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

func ParsePostUserRequest(r *http.Request) (user *db.User, err error) {
	var newUser NewUserRequest
	var uuidStr string
	var uuid *gocql.UUID
	uuid = nil
	uuidStr = r.Header.Get("X-UUID")
	fmt.Printf("X-UUID: '%v' \n", uuidStr)
	if uuidStr != "" {
		tmpuuid, err := gocql.ParseUUID(uuidStr)
		if err != nil {
			fmt.Printf("Failed to parse X-UUID: %v\n", uuidStr)
		} else {
			uuid = &tmpuuid
		}
	}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&newUser)
	if err != nil {
		fmt.Println("Parse Post User Request errors", err.Error())
		return nil, err
	}
	if newUser.Name == nil || len(*newUser.Name) == 0 {
		return nil, errors.New("User: Name is missing or empty")
	}
	return &db.User{Name: *newUser.Name, Id: uuid}, err
}

func PostUser(w http.ResponseWriter, r *http.Request) {
	var user *db.User
	var err error
	user, err = ParsePostUserRequest(r)
	if err != nil {
		HandleHttpResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	uuid, err := DB.AddUser(*user)
	if err != nil {
		HandleHttpResponse(w, http.StatusInternalServerError, "Failed to persist "+user.Id.String()+": "+err.Error())
		return
	}
	user.Id = uuid
	HandleNewUserResponse(w, *user)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	var uuid, error = ParseUserid(r)
	if error != nil {
		HandleHttpResponse(w, http.StatusBadRequest, error.Error())
		return
	}
	user, error := DB.GetUser(uuid)
	if error != nil {
		HandleHttpResponse(w, http.StatusNotFound, error.Error())
		return
	}
	HandleUserResponse(w, *user)
	return
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var uuid, error = ParseUserid(r)
	if error != nil {
		fmt.Println("Failed to parse user id")
		HandleHttpResponse(w, http.StatusBadRequest, error.Error())
		return
	}
	fmt.Printf("userid %v", uuid)
	error = DB.DeleteUser(uuid)
	if error != nil {
		fmt.Printf("Failed to delete user %v: %v\n", uuid, error.Error())
		// Handle not found and internal server error
		HandleHttpResponse(w, http.StatusInternalServerError, error.Error())
	}

}

func HandleUserResponse(w http.ResponseWriter, user db.User) {
	fmt.Println("user_id", user.Id)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "   ")
	encoder.Encode(user)
}

func HandleNewUserResponse(w http.ResponseWriter, user db.User) {
	fmt.Println("user_id", user.Id)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "   ")
	encoder.Encode(UserResponse{Id: *user.Id, Name: user.Name})
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	var users []UserResponse
	userList, err := DB.ListUsers()
	if err != nil {
		HandleHttpResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, user := range userList {
		users = append(users, UserResponse{
			Id:   *user.Id,
			Name: user.Name,
		})
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "   ")
	encoder.Encode(UsersResponse{Users: users})
}

func HandleErrorsResponse(w http.ResponseWriter, error error) {
	fmt.Println("error: ", error)
	json.NewEncoder(w).Encode(ErrorResponse{Error: error})
}
