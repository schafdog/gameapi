package main

import (
	"github.com/gorilla/mux"
	"github.com/schafdog/gameapi/db"
	"github.com/schafdog/gameapi/user"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     db.UserDatabase
}

func (a *App) Initialize(newDB db.UserDatabase) {
	var err error

	if err != nil {
		panic(err)
	}
	a.DB = newDB
	User.InitDB(newDB)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8000", a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.StrictSlash(true)
	a.Router.HandleFunc("/user", User.PostUser).Methods("POST")
	a.Router.HandleFunc("/user", User.GetUsers).Methods("GET")
	a.Router.HandleFunc("/user/{userid}", User.DeleteUser).Methods("DELETE")
	a.Router.HandleFunc("/user/{userid}/state", User.GetState).Methods("GET")
	a.Router.HandleFunc("/user/{userid}/state", User.PutState).Methods("PUT")
	a.Router.HandleFunc("/user/{userid}/friends", User.GetFriends).Methods("GET")
	a.Router.HandleFunc("/user/{userid}/friends", User.PutFriends).Methods("PUT")
	a.Router.HandleFunc("/user/{userid}", User.GetUser).Methods("GET")

	/*
		a.Router.HandleFunc("/user/{userid}", User.Put).Methods("PUT")
		a.Router.HandleFunc("/user/{userid}", User.Patch).Methods("PATCH")
		a.Router.HandleFunc("/user/{UUID}", User.Patch).Methods("HIGHSCORE")
	*/

}

func (a *App) Close() {
	a.DB.Close()
}
