package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"."
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))
	//	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if err := a.Session.Query(tableCreationQuery).Exec(); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.Session.Query("TRUNCATE user").Exec()
}

const tableCreationQuery = `
use gameapi;
drop table user;
create table user (
  id UUID,
  name text,
  score int,
  gamesplayed int,
  friends set<UUID>,
  version int,
  PRIMARY KEY(id)
);`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/user", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	/*
		if body := response.Body.String(); body != "{ [] }" {
			t.Errorf("Expected an empty array. Got %s", body)
		}
	*/
}

func TestGetNonExistentUser(t *testing.T) {
	var userid gocql.UUID
	userid = gocql.TimeUUID()
	clearTable()
	req, _ := http.NewRequest("GET", "/user/"+userid.String(), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)

	/*
		var m map[string]string
		json.Unmarshal(response.Body.Bytes(), &m)
		if m["error"] != "User not found" {
			t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
		}
	*/
}

// main_test.go

func createUser(name string, t *testing.T) string {
	payload := []byte(fmt.Sprintf(`{ "name": "%v" }`, name))

	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	fmt.Printf("Create User Response %v\n", m)
	fmt.Printf("Id %v\n", m["id"])

	_, error := gocql.ParseUUID(m["id"].(string))
	if error != nil {
		t.Errorf("Expected user ID to be UUID. Got '%v'", m["id"])
	}
	fmt.Printf("NAME %v\n", m["name"])

	if m["name"] != name {
		t.Errorf("Expected user name to be '%v'. Got '%v'", name, m["name"])
	}
	return m["id"].(string)
}

func TestCreateUser(t *testing.T) {
	clearTable()

	createUser("test user", t)

}

func setGetState(score int, gamesPlayed int, userid string, t *testing.T) {
	payload := []byte(fmt.Sprintf(`{ "score": %v, "gamesPlayed": %v }`, score, gamesPlayed))
	path := "/user/" + userid + "/state"
	req, _ := http.NewRequest("PUT", path, bytes.NewBuffer(payload))
	response1 := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response1.Code)
	fmt.Printf("PUT state %v \n", response1.Code)

	req, _ = http.NewRequest("GET", path, nil)
	response2 := executeRequest(req)
	fmt.Printf("GET response: %v \n", response2.Code)
	checkResponseCode(t, http.StatusOK, response2.Code)
	var m map[string]interface{}
	json.Unmarshal(response2.Body.Bytes(), &m)

	fmt.Printf("GET state %v\n", m)

	if int(m["score"].(float64)) != score {
		t.Errorf("Expected score to be '%v'. Got '%v'", score, m["score"])
	}

	if int(m["gamesPlayed"].(float64)) != gamesPlayed {
		t.Errorf("Expected games played to be '%v'. Got '%v'", gamesPlayed, m["gamesPlayed"])
	}
}

func TestResetUserState(t *testing.T) {
	clearTable()
	userid := createUser("ResetStatUser", t)
	setGetState(0, 0, userid, t)
}

func TestSetUserState(t *testing.T) {
	clearTable()
	userid := createUser("SetStatUser", t)
	setGetState(100, 9, userid, t)
}
