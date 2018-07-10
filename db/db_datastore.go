package db_model

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/gocql/gocql"
	"golang.org/x/net/context"
)

// datastoreDB persists users to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type datastoreDB struct {
	client *datastore.Client
}

// Ensure datastoreDB conforms to the UserDatabase interface.
var _ UserDatabase = &datastoreDB{}

// newDatastoreDB creates a new UserDatabase backed by Cloud Datastore.
// See the datastore and google packages for details on creating a suitable Client:
// https://godoc.org/cloud.google.com/go/datastore
func newDatastoreDB(client *datastore.Client) (UserDatabase, error) {
	ctx := context.Background()
	// Verify that we can communicate and authenticate with the datastore service.
	t, err := client.NewTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	if err := t.Rollback(); err != nil {
		return nil, fmt.Errorf("datastoredb: could not connect: %v", err)
	}
	return &datastoreDB{
		client: client,
	}, nil
}

// Close closes the database.
func (db *datastoreDB) Close() error {
	// No op.
	return nil
}

func (db *datastoreDB) datastoreKey(userid gocql.UUID) *datastore.Key {

	return datastore.IDKey(userid.String(), 0, nil)
}

// GetUser retrieves a user by its ID.
func (db *datastoreDB) GetUser(userid gocql.UUID) (*User, error) {
	ctx := context.Background()
	k := db.datastoreKey(userid)
	user := &User{}
	if err := db.client.Get(ctx, k, user); err != nil {
		return nil, fmt.Errorf("datastoredb: could not get User: %v", err)
	}
	user.Id = &userid
	return user, nil
}

// AddUser saves a given user, assigning it a new ID.
func (db *datastoreDB) AddUser(user User) (userid *gocql.UUID, err error) {
	ctx := context.Background()
	id := gocql.TimeUUID()
	k := datastore.IncompleteKey("User", nil)
	k, err = db.client.Put(ctx, k, user)
	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not put User: %v", err)
	}
	return &id, fmt.Errorf("Not Implemented: Returning wrong key")
}

// DeleteUser removes a given user by its ID.
func (db *datastoreDB) DeleteUser(userid gocql.UUID) error {
	ctx := context.Background()
	k := db.datastoreKey(userid)
	if err := db.client.Delete(ctx, k); err != nil {
		return fmt.Errorf("datastoredb: could not delete User: %v", err)
	}
	return nil
}

// UpdateUser updates the entry for a given user.
func (db *datastoreDB) UpdateUser(user User) error {
	ctx := context.Background()
	k := db.datastoreKey(*user.Id)
	if _, err := db.client.Put(ctx, k, user); err != nil {
		return fmt.Errorf("datastoredb: could not update User: %v", err)
	}
	return nil
}

// ListUsers returns a list of users, ordered by title.
func (db *datastoreDB) ListUsers() ([]User, error) {
	ctx := context.Background()
	users := make([]User, 0)
	q := datastore.NewQuery("gameapi").Order("name")

	_, err := db.client.GetAll(ctx, q, &users)

	if err != nil {
		return nil, fmt.Errorf("datastoredb: could not list users: %v", err)
	}
	// Not sure how to get UUID
	/*
		for i, k := range keys {
			users[i].Id = gocql.GenerateUUID(k.ID)
		}
	*/
	return users, nil
}

// SetState sets the state of a User
func (db *datastoreDB) SetState(userid gocql.UUID, state State) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *datastoreDB) GetState(userid gocql.UUID) (*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// SetFriends sets friends of the user
func (db *datastoreDB) SetFriends(userid gocql.UUID, friends []gocql.UUID) error {
	return fmt.Errorf("Not Implemeted")
}

// GetState returns the current state of a User
func (db *datastoreDB) GetFriends(userid gocql.UUID) ([]*gocql.UUID, error) {
	return nil, fmt.Errorf("Not Implemeted")
}

// GetFriendsState returns a list of users, ordered by title, filtered by
// the user who created the user entry.
func (db *datastoreDB) GetFriendsState(userid gocql.UUID) ([]*State, error) {
	return nil, fmt.Errorf("Not Implemeted")
}
