package db_model

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/gocql/gocql"
)

// datastoreDB persists users to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type templateDB struct {
}

// Ensure templateDB conforms to the UserDatabase interface.
var _ UserDatabase = &templateDB{}

// newDB is a template for new DBs
func newDB(dbUrl string) (UserDatabase, error) {
	return &templateDB{}, nil
}

// Close closes the database.
func (db *templateDB) Close() error {
	// No op.
	return fmt.Errorf("Not Implemented")
}

func (db *templateDB) datastoreKey(userid gocql.UUID) *datastore.Key {
	return nil
}

// GetUser retrieves a user by its ID.
func (db *templateDB) GetUser(userid gocql.UUID) (*User, error) {
	return nil, fmt.Errorf("Not implemented")
}

// AddUser saves a given user, assigning it a new ID.
func (db *templateDB) AddUser(user User) (userid *gocql.UUID, err error) {
	return nil, fmt.Errorf("Not Implemented: Returning wrong key")
}

// DeleteUser removes a given user by its ID.
func (db *templateDB) DeleteUser(userid gocql.UUID) error {
	return fmt.Errorf("Not Implemented")
}

// UpdateUser updates the entry for a given user.
func (db *templateDB) UpdateUser(user User) error {
	return fmt.Errorf("Not Implemented")
}

// ListUsers returns a list of users, ordered by title.
func (db *templateDB) ListUsers() ([]User, error) {
	return []User{}, fmt.Errorf("Not Implemented")
}

// SetState sets the state of a User
func (db *templateDB) SetState(userid gocql.UUID, state State) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *templateDB) GetState(userid gocql.UUID) (*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// SetFriends sets friends of the user
func (db *templateDB) SetFriends(userid gocql.UUID, friends []gocql.UUID) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *templateDB) GetFriends(userid gocql.UUID) ([]*gocql.UUID, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// GetFriendsState returns a list of users, ordered by title, filtered by
// the user who created the user entry.
func (db *templateDB) GetFriendsState(userid gocql.UUID) ([]*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}
