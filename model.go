package main

import (
	//	"errors"
	"github.com/gocql/gocql"
)

type User struct {
	Id          *gocql.UUID  `json:"id"`
	Name        string       `json:"name"`
	Highscore   int          `json:"score"`
	GamesPlayed int          `json:"gamesPlayed"`
	Friends     []gocql.UUID `json:"friends"`
}

// State to form payload returning a single User state
type State struct {
	Id          gocql.UUID `json:"id"`
	Highscore   int        `json:"score"`
	GamesPlayed int        `json:"gamesPlayed"`
}

type UserDatabase interface {

	// CRUD operations for the Model

	// Create a User
	CreateUser(user User) (uuid *gocql.UUID, err error)
	RetrieveUser(userid gocql.UUID) (user *User, err error)
	UpdateUser(user User) error

	// Delete a user
	DeleteUser(useri gocql.UUID) (err error)

	ListUsers() (users []User, err error)

	// Set, Get for User state
	SetState(userid gocql.UUID, state State) error
	GetState(userid gocql.UUID) (state *State, err error)

	// Set, Get for User friends
	SetFriends(userid gocql.UUID, friends []gocql.UUID) error
	GetFriendsState(userid gocql.UUID) (state []State, err error)

	// Close any outstanding DB resources
	Close() error
}
