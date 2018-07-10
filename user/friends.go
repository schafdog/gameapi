package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"net/http"
)

type FriendsRequest struct {
	Friends []gocql.UUID `json:"friends"`
}

// Friends builds a payload of new user resource ID
type FriendResponse struct {
	Id        gocql.UUID `json:"id"`
	Name      string     `json:"name"`
	Highscore int        `json:"highscore"`
}

type FriendsResponse struct {
	Friends []FriendResponse `json:"friends"`
}

func getFriends(userid gocql.UUID) (FriendsResponse, error) {
	var friends []gocql.UUID
	var friendsList []FriendResponse
	friendsState, err := DB.GetFriendsState(userid)
	return FriendsResponse{Friends: friendsState}, err
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	var uuid, error = ParseUserid(r)
	if error != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: error})
		return
	}
	friends, err := getFriends(uuid)
	handleFriendsResponse(w, friends, err)
}

func ParseFriendsRequest(r *http.Request) (FriendsRequest, error) {
	var friendsRequest FriendsRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&friendsRequest)
	fmt.Printf("New friends: %v \n", friendsRequest)
	if err != nil {
		fmt.Println("Error parse FriendsRequest ", err.Error())
	}
	return friendsRequest, err

}

func PutFriends(w http.ResponseWriter, r *http.Request) {
	var friendsRequest FriendsRequest
	var uuid, error = ParseUserid(r)
	if error != nil {
		json.NewEncoder(w).Encode(ErrorResponse{Error: error})
		return
	}
	fmt.Printf("PUT Friends: User Id: %v\n", uuid)
	friendsRequest, error = ParseFriendsRequest(r)
	if error != nil {
		HandleHttpResponse(w, http.StatusBadRequest, error.Error())
		return
	}
	fmt.Printf("PUT Friends: UpdateFriends: %v %v\n", uuid, friendsRequest.Friends)
	err := DB.UpdateFriends(uuid, friendsRequest)
	var status = http.StatusOK
	var message = ""
	if error != nil {
		message = error.Error()
		if message == "Not Found" {
			status = http.StatusNotFound
		} else {
			status = http.StatusBadRequest
		}
	}
	HandleHttpResponse(w, status, message)
}

func handleFriendsResponse(w http.ResponseWriter, friends FriendsResponse, err error) {
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "   ")
		encoder.Encode(friends)
	} else {
		fmt.Println("errors", err)
		json.NewEncoder(w).Encode(ErrorResponse{Error: err})
	}
}

func HandleHttpResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if message != "" || len(message) > 0 {
		json.NewEncoder(w).Encode(message)
	}
}
