package main

import (
	"database/sql"
)

type post struct {
	ID     int    `json:"ID"`
	Title  string `json:"title"`
	Body   string `json:"body"`
	Author string `json:"username"`
}

func (p *post) createPost(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO posts(title, body, author) VALUES($1, $2, $3) RETURNING id",
		p.Title, p.Body, p.Author,
	).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

func (p *post) getPost(db *sql.DB) error {
	return db.QueryRow(
		"SELECT title, body, author FROM posts WHERE id=$1", p.ID,
	).Scan(&p.Title, &p.Body, &p.Author)
}

func (p *post) updatePost(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE posts SET title=$1, body=$2, author=$3 WHERE id=$4",
			p.Title, p.Body, p.Author, p.ID)

	return err
}

func (p *post) deletePost(db *sql.DB) error {
	_, err :=
		db.Exec("DELETE FROM posts WHERE id=$1", p.ID)

	return err
}

func getAllPosts(db *sql.DB) ([]post, error) {
	rows, err := db.Query(
		"SELECT * FROM posts;",
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []post{}

	for rows.Next() {
		var p post
		if err := rows.Scan(&p.ID, &p.Title, &p.Body, &p.Author); err != nil {
			return nil, err
		}
	}

	return posts, err
}
