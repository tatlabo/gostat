package models

import (
	"database/sql"
	"time"
)

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (*int, error) {

	stmt := `INSERT INTO snippets (title, content, created, expires) VALUES ($1, $2, $3, $4) RETURNING id;`
	var id int
	err := m.DB.QueryRow(stmt, s.Title, s.Content, s.Created, s.Expires).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	return nil, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
