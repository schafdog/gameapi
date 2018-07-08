package User

import (
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/schafdog/gameapi/cassandra"
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

func getFriends(userid gocql.UUID) (FriendsResponse, []string) {
	var friends []gocql.UUID
	var friendsList []FriendResponse
	var errs []string
	if err := cassandra.Session.Query(`
      select friends2 from User where id = ?`,
		userid).Scan(&friends); err != nil {
		errs = append(errs, err.Error())
		return FriendsResponse{Friends: friendsList}, errs
	}
	fmt.Printf("getFriends found: %v\n", friends)

	m := map[string]interface{}{}
	query := "select id, name, score from user where id in ?"
	iterable := cassandra.Session.Query(query, friends).Iter()
	if iterable == nil {
		fmt.Printf("Failed to iter\n")
	}
	for iterable.MapScan(m) {
		fmt.Printf("User{ id: %v, name: %v, highscore: %v }\n", m["id"], m["name"], m["score"])
		friendsList = append(friendsList, FriendResponse{
			Id:        m["id"].(gocql.UUID),
			Name:      m["name"].(string),
			Highscore: m["score"].(int),
		})
		m = map[string]interface{}{}
	}
	return FriendsResponse{Friends: friendsList}, errs
}

func GetFriends(w http.ResponseWriter, r *http.Request) {
	var errs []string
	var uuid, error = ParseUserid(r)
	if error != nil {
		errs = append(errs, "Failed to parse uuid: "+error.Error())
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
		return
	}
	friends, errs := getFriends(uuid)
	handleFriendsResponse(w, friends, errs)
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

func UpdateFriends(userid gocql.UUID, friends FriendsRequest) error {
	if err := cassandra.Session.Query(`update User set friends2 = ? where id = ?`,
		friends.Friends, userid).Exec(); err != nil {
		fmt.Println("Failed to update friends: ", err.Error())
		return err
	}
	return nil
}

func PutFriends(w http.ResponseWriter, r *http.Request) {
	var friendsRequest FriendsRequest
	var uuid, error = ParseUserid(r)
	if error != nil {
		var errs []string
		errs = append(errs, "Failed to parse uuid: "+error.Error())
		json.NewEncoder(w).Encode(ErrorResponse{Errors: []string{error.Error()}})
		return
	}
	fmt.Printf("PUT Friends: User Id: %v\n", uuid)
	friendsRequest, error = ParseFriendsRequest(r)
	if error != nil {
		HandleHttpResponse(w, http.StatusBadRequest, error.Error())
		return
	}
	fmt.Printf("PUT Friends: UpdateFriends: %v %v\n", uuid, friendsRequest.Friends)
	error = UpdateFriends(uuid, friendsRequest)
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

func handleFriendsResponse(w http.ResponseWriter, friends FriendsResponse, errs []string) {
	if errs == nil {
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "   ")
		encoder.Encode(friends)
	} else {
		fmt.Println("errors", errs)
		json.NewEncoder(w).Encode(ErrorResponse{Errors: errs})
	}
}

func HandleHttpResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	if message != "" || len(message) > 0 {
		json.NewEncoder(w).Encode(message)
	}

	/*
		if errs == nil {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "   ")
			encoder.Encode(string)
		} else {
			fmt.Println("errors", errs)
			json.NewEncoder(w).Encode(ErrorResponse{Errors: message})
		}
	*/
}
