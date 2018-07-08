package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/schafdog/gameapi/cassandra"
	"github.com/schafdog/gameapi/user"
	"log"
	"net/http"
)

type userCreateResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userCreateResponse{ID: "UUID", Name: "John"})
}

func main() {
	CassandraSession := cassandra.Session
	defer CassandraSession.Close()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/user", User.PostUser).Methods("POST")
	router.HandleFunc("/user", User.GetUsers).Methods("GET")
	router.HandleFunc("/user/{userid}/state", User.GetStat).Methods("GET")
	router.HandleFunc("/user/{userid}/state", User.PutStat).Methods("PUT")
	router.HandleFunc("/user/{userid}/friends", User.GetFriends).Methods("GET")
	router.HandleFunc("/user/{userid}/friends", User.PutFriends).Methods("PUT")
	/*
		router.HandleFunc("/user/{userid}", User.Delete).Methods("DELETE")
		   router.HandleFunc("/user/{userid}", User.Get).Methods("GET")
		   router.HandleFunc("/user/{userid}", User.Put).Methods("PUT")
		   router.HandleFunc("/user/{userid}", User.Patch).Methods("PATCH")
		   router.HandleFunc("/user/{UUID}", User.Patch).Methods("HIGHSCORE")
	*/
	log.Fatal(http.ListenAndServe(":8000", router))
}
