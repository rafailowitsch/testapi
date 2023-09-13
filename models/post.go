package main

import (
	"database/sql"
	"errors"
)

type post struct {
	title  string `json:"title"`
	body   string `json:"body"`
	author *user
}

func (p *post) createPost(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *post) getPost(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *post) updatePost(db *sql.DB) error {
	return errors.New("Not implemented")
}

func (p *post) deletePost(db *sql.DB) error {
	return errors.New("Not implemented")
}

func getAllPosts(db *sql.DB) ([]post, error) {
	return nil, errors.New("Not implemented")
}
