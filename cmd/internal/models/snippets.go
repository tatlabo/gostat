package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Snippet struct {
	ID      int       `db:"id" json:"id"`
	Title   string    `db:"title" json:"title"`
	Content string    `db:"content" json:"content"`
	Created time.Time `db:"created" json:"created"`
	Expires string    `db:"expires" json:"expires"`
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (*int, error) {

	stmt := `INSERT INTO snippets (title, content, expires, created) 
		VALUES ($1, $2, to_date($3, 'YYYY-MM-DD'), NOW()) RETURNING id;`

	fmt.Printf("\n%#v\n", s)

	var id int
	err := m.DB.QueryRow(stmt, s.Title, s.Content, s.Expires).Scan(&id)
	fmt.Printf("stmt:\n %s, %s, %s, %v\n", stmt, s.Title, s.Content, s.Expires)
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
