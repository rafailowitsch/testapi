package main

import (
	"database/sql"
	"errors"
)

type user struct {
	username string `json:"username"`
	password string `json:"password"`
}

func (u *user) createUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (u *user) deleteUser(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getAllUsers(db *sql.DB) ([]user, error) {
	return nil, errors.New("Not implemented")
}
