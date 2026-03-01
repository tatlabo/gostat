package models

import (
	"database/sql"
	"errors"
	"fmt"
	"gostats/cmd/internal/highlight"
	"html/template"
	"time"
)

type Snippet struct {
	ID      int           `db:"id" json:"id"`
	Title   string        `db:"title" json:"title"`
	Content string        `db:"content" json:"content"`
	Created time.Time     `db:"created" json:"created"`
	Expires time.Time     `db:"expires" json:"expires"`
	Html    template.HTML `db:"-" json:"html"`
}

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(s *Snippet) (*int, error) {

	stmt := `INSERT INTO snippets (title, content, expires, created) 
		VALUES ($1, $2, to_date($3, 'YYYY-MM-DD'), NOW()) RETURNING id;`

	var id int
	err := m.DB.QueryRow(stmt, s.Title, s.Content, s.Expires).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &id, nil
}

func (m *SnippetModel) Delete(s *Snippet) (Snippet, error) {

	stmt := `--sql
		DELETE FROM snippets WHERE id = $1 RETURNING title, content, expires, created;`

	err := m.DB.QueryRow(stmt, s.ID).Scan(&s.Title, &s.Content, &s.Expires, &s.Created)

	if err != nil {
		return Snippet{}, err
	}

	return *s, nil
}

func (s *Snippet) HighlightSnippet() (err error) {

	html, err := highlight.Highlight(s.Content)
	s.Html = template.HTML(html)

	if err != nil {
		return fmt.Errorf("error highlighting snippet content: %v", err)
	}

	return nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {

	stmt := `SELECT id, title, content, expires, created FROM snippets
		WHERE id = $1;`

	s := &Snippet{}
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, err

}

func (m *SnippetModel) Latest() (list []Snippet, err error) {

	stmt := `SELECT id, title, content, expires, created FROM snippets ORDER BY created DESC LIMIT 100;`

	res, err := m.DB.Query(stmt)
	if err != nil {
		return nil,
			fmt.Errorf("Error executing query: SELECT id, title, content...: %v", err)
	}
	defer res.Close()

	for res.Next() {
		s := Snippet{}
		res.Scan(&s.ID, &s.Title, &s.Content, &s.Expires, &s.Created)
		if err != nil {
			return nil, err
		}
		list = append(list, s)

	}

	return list, err

}
