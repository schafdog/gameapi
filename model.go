package main

import (
	"database/sql"
	"errors"
	"github.com/gocql/gocql"
)

type user struct {
	Id          gocql.UUID   `json:"id"`
	Name        string       `json:"name"`
	Highscore   int          `json:"score"`
	GamesPlayed int          `json:"gamesPlayed"`
	Friends     []gocql.UUID `json:"friends"`
}

func (p *user) getUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *user) updateUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *user) deleteUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *user) createUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getUsers(db *sql.DB, start, count int) ([]user, error) {
	return nil, errors.New("Not implemented")
}
