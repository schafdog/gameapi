package main

import (
	"github.com/schafdog/gameapi/db"
)

var a *App

func main() {
	a := App{}
	DB, error := db.NewCassandraDB("127.0.0.1;gameapi;")
	if error != nil {
		panic("Error opening DB: " + error.Error())
	}
	a.Initialize(DB)

	defer DB.Close()
	a.Run(":8000")
}
