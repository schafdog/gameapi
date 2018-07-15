package db

import (
	"cloud.google.com/go/datastore"
	"fmt"
	"github.com/gocql/gocql"
)

// datastoreDB persists users to Cloud Datastore.
// https://cloud.google.com/datastore/docs/concepts/overview
type CassandraDB struct {
	Session *gocql.Session
}

// Ensure CassandraDB conforms to the UserDatabase interface.
var _ UserDatabase = &CassandraDB{}

// newDB is a template for new DBs
func NewCassandraDB(dbUrl string) (db UserDatabase, err error) {
	this := CassandraDB{}
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "gameapi"
	this.Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("cassandra init done")

	return &this, nil
}

// Close closes the database.
func (db *CassandraDB) Close() error {
	// No op.
	fmt.Printf("Closing Cassandra DB")
	db.Session.Close()
	return nil
}

func (db *CassandraDB) datastoreKey(userid gocql.UUID) *datastore.Key {
	return nil
}

// GetUser retrieves a user by its ID.
func (db *CassandraDB) GetUser(userid gocql.UUID) (user *User, err error) {
	user = &User{}
	query := "SELECT id,name, score, gamesPlayed  FROM User where id = ?"
	if err := db.Session.Query(query, userid).
		Scan(&user.Id, &user.Name, &user.GamesPlayed, &user.Highscore); err != nil {
		fmt.Printf("Failed to find user for %v", userid)
		return nil, err
	}
	fmt.Printf("Found user %v for %v", user.Name, userid)
	return user, nil
}

// AddUser saves a given user, assigning it a new ID.
func (db *CassandraDB) AddUser(user User) (userid *gocql.UUID, err error) {
	var gocqlUuid gocql.UUID

	// generate a unique UUID for this user
	if user.Id == nil {
		gocqlUuid = gocql.TimeUUID()
		fmt.Printf("creating a new userid %v for %v\n", gocqlUuid, user.Name)
		user.Id = &gocqlUuid
	} else {
		fmt.Printf("Using suggestion %v from header for %v\n", user.Id, user.Name)
	}
	// write data to Cassandra
	var friends []gocql.UUID
	err = db.Session.Query(`INSERT INTO user (id, name, score, gamesplayed, friends) VALUES (?, ?, ?, ?, ?)`, user.Id, user.Name, 0, 0, friends).Exec()
	return user.Id, err
}

// DeleteUser removes a given user by its ID.
func (db *CassandraDB) DeleteUser(userid gocql.UUID) error {
	// write data to Cassandra
	err := db.Session.Query(`DELETE FROM user where id = ?`, userid).Exec()
	return err
}

// UpdateUser updates the entry for a given user.
func (db *CassandraDB) UpdateUser(user User) error {
	return fmt.Errorf("Not Implemented")
}

// ListUsers returns a list of users, ordered by title.
func (db *CassandraDB) ListUsers() ([]User, error) {
	var userList []User
	query := "SELECT id,name FROM User"
	m := map[string]interface{}{}
	iterable := db.Session.Query(query).Iter()
	for iterable.MapScan(m) {
		fmt.Printf("User{ id: %v, name: %v }\n", m["id"], m["name"])
		uuid := m["id"].(gocql.UUID)
		userList = append(userList, User{
			Id:   &uuid,
			Name: m["name"].(string),
		})
		m = map[string]interface{}{}
	}
	return userList, nil
}

// SetState sets the state of a User
func (db *CassandraDB) SetState(userid gocql.UUID, state State) error {
	if err := db.Session.Query(`
      UPDATE user set gamesPlayed = ?, score = ? where id = ?`,
		state.GamesPlayed, state.Highscore, state.Id).Exec(); err != nil {
	}
	return nil
}

// GetState returns the current state of a User
func (db *CassandraDB) GetState(userid gocql.UUID) (*State, error) {
	state := State{Id: userid}
	err := db.Session.Query(`
      select gamesPlayed, score from User where id = ?`,
		state.Id).Scan(&state.GamesPlayed, &state.Highscore)
	return &state, err
}

// SetFriends sets friends of the user
func (db *CassandraDB) SetFriends(userid gocql.UUID, friends []gocql.UUID) error {
	if err := db.Session.Query(`update User set friends = ? where id = ?`,
		friends, userid).Exec(); err != nil {
		fmt.Println("Failed to update friends: ", err.Error())
		return err
	}
	return nil
}

// GetFriends returns the friends of a User
// Not public API
func (db *CassandraDB) GetFriends(userid gocql.UUID) (friendsList []*gocql.UUID, err error) {
	var friends []*gocql.UUID
	if err := db.Session.Query(`
      select friends from User where id = ?`,
		userid).Scan(&friends); err != nil {
		return nil, fmt.Errorf("Failed to get friends from %v: %v", userid, err.Error())
	}
	fmt.Printf("getFriends found: %v\n", friends)
	return friends, nil
}

// GetFriendsState returns a list of users, ordered by title, filtered by
// the user who created the user entry.
func (db *CassandraDB) GetFriendsState(userid gocql.UUID) (friendsState []*State, err error) {
	friends, err := db.GetFriends(userid)
	if err != nil {
		return nil, err
	}
	query := "select id, name, score from user where id in ?"
	iterable := db.Session.Query(query, friends).Iter()
	if iterable == nil {
		fmt.Printf("Failed to iter\n")
	}
	m := map[string]interface{}{}
	for iterable.MapScan(m) {
		fmt.Printf("User{ id: %v, name: %v, highscore: %v }\n", m["id"], m["name"], m["score"])
		friendsState = append(friendsState, &State{
			Id:        m["id"].(gocql.UUID),
			Name:      m["name"].(string),
			Highscore: m["score"].(int),
		})
		m = map[string]interface{}{}
	}
	return friendsState, nil
}
