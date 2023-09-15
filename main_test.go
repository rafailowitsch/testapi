package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
	)

	ensureTableExists()
	code := m.Run()
	clearTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tablesCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM posts")
	a.DB.Exec("ALTER SEQUENCE posts_id_seq RESTART WITH 1")
}

const tablesCreationQuery = `
CREATE TABLE posts (
	id SERIAL NOT NULL UNIQUE,
	title VARCHAR(255) NOT NULL,
	body VARCHAR(255) NOT NULL,
	author VARCHAR(255) NOT NULL
);
`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, err := http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatalf("Could not create request for post: %v", err)
	}

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)

	}

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder() // new responserecorder
	a.Router.ServeHTTP(rr, req)  // processes the request and writes it to the responserecorder

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentPost(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "post/999999", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Post not found" {
		err := fmt.Sprintf("Expected the 'error' key of the respose to be set to 'Post not found'. Got '%s'.",
			m["error"])
		t.Errorf(err)
	}
}

func TestCreatePost(t *testing.T) {
	clearTable()

	var jsonStr = []byte(
		`{"title":"why do I feel bad?",
		  "body":"because you're weak.",
		  "author":"Einstein",}`,
	)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["title"] != "why do I feel bad?" {
		t.Errorf("Expected post title to be 'why do I feel bad?'. Got '%v'", m["title"])
	}
	if m["body"] != "because you're weak." {
		t.Errorf("Expected post body to be 'because you're weak.'. Got '%v'", m["body"])
	}
	if m["author"] != "Einstein" {
		t.Errorf("Expected post body to be 'Einstein'. Got '%v'", m["author"])
	}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetPost(t *testing.T) {
	clearTable()
	addPosts(1)

	req, _ := http.NewRequest("GET", "/post/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdatePost(t *testing.T) {
	clearTable()
	addPosts(1)

	req, _ := http.NewRequest("GET", "/post/1", nil)
	response := executeRequest(req)
	var originalPost map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalPost)

	var jsonStr = []byte(
		`{"title":"test post",
		  "body":"test body",
		  "author":"test Einstein",}`,
	)
	req, _ = http.NewRequest("PUT", "/product", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalPost["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v.",
			originalPost["id"], m["id"])
	}

	if m["title"] != originalPost["title"] {
		t.Errorf("Expected the title to change from %v to %v. Got %v.",
			originalPost["title"], m["title"], m["title"])
	}

	if m["body"] != originalPost["body"] {
		t.Errorf("Expected the body to change from %v to %v. Got %v.",
			originalPost["body"], m["body"], m["body"])
	}

	if m["author"] != originalPost["author"] {
		t.Errorf("Expected the author to change from %v to %v. Got %v.",
			originalPost["author"], m["author"], m["author"])
	}
}

func TestDeletePost(t *testing.T) {
	clearTable()
	addPosts(1)

	req, _ := http.NewRequest("GET", "/post/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/post/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/post/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// add test post(-s)
func addPosts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO posts(title, body, author) VALUES($1, $2, $3)",
			"Post "+strconv.Itoa(i), "Body "+strconv.Itoa(i), "Author "+strconv.Itoa(i))
	}
}
