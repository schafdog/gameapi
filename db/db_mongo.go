package db_model

import (
	"crypto/rand"
	"fmt"
	"github.com/gocql/gocql"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math/big"
)

type mongoDB struct {
	conn *mgo.Session
	c    *mgo.Collection
}

// Ensure mongoDB conforms to the UserDatabase interface.
var _ UserDatabase = &mongoDB{}

// newMongoDB creates a new UserDatabase backed by a given Mongo server,
// authenticated with given credentials.
func newMongoDB(addr string, cred *mgo.Credential) (UserDatabase, error) {

	return nil, fmt.Errorf("Not Implement")

	conn, err := mgo.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("mongo: could not dial: %v", err)
	}

	if cred != nil {
		if err := conn.Login(cred); err != nil {
			return nil, err
		}
	}

	return &mongoDB{
		conn: conn,
		c:    conn.DB("gameapi").C("users"),
	}, nil
}

// Close closes the database.
func (db *mongoDB) Close() error {
	db.conn.Close()
	return nil
}

// GetUser retrieves a user by its ID.
func (db *mongoDB) GetUser(userid gocql.UUID) (*User, error) {
	user := &User{}
	if err := db.c.Find(bson.D{{Name: "id", Value: userid}}).One(user); err != nil {
		return nil, err
	}
	return user, nil
}

var maxRand = big.NewInt(1<<63 - 1)

// randomID returns a positive number that fits within an int64.
func randomID() (int64, error) {
	// Get a random number within the range [0, 1<<63-1)
	n, err := rand.Int(rand.Reader, maxRand)
	if err != nil {
		return 0, err
	}
	// Don't assign 0.
	return n.Int64() + 1, nil
}

// AddUser saves a given user, assigning it a new ID.
func (db *mongoDB) AddUser(user User) (userid *gocql.UUID, err error) {

	var id = gocql.TimeUUID()
	if err != nil {
		return nil, fmt.Errorf("mongodb: could not assign an new ID: %v", err)
	}

	user.Id = &id
	if err := db.c.Insert(user); err != nil {
		return nil, fmt.Errorf("mongodb: could not add user: %v", err)
	}
	return &id, nil
}

// DeleteUser removes a given user by its ID.
func (db *mongoDB) DeleteUser(userid gocql.UUID) error {
	return db.c.Remove(bson.D{{Name: "id", Value: userid}})
}

// UpdateUser updates the entry for a given user.
func (db *mongoDB) UpdateUser(user User) error {
	return db.c.Update(bson.D{{Name: "id", Value: user.Id}}, user)
}

// ListUsers returns a list of users, ordered by title.
func (db *mongoDB) ListUsers() ([]User, error) {
	var result []User
	if err := db.c.Find(nil).Sort("name").All(&result); err != nil {
		return nil, err
	}
	return result, nil
}

// SetState sets the state of a User
func (db *mongoDB) SetState(userid gocql.UUID, state State) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *mongoDB) GetState(userid gocql.UUID) (*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// SetFriends sets friends of the user
func (db *mongoDB) SetFriends(userid gocql.UUID, friends []gocql.UUID) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *mongoDB) GetFriends(userid gocql.UUID) ([]*gocql.UUID, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// GetFriendsState returns a list of users, ordered by title, filtered by
// the user who created the user entry.
func (db *mongoDB) GetFriendsState(userid gocql.UUID) ([]*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}
