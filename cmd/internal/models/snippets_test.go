package models

import (
	"database/sql"
	"gostats/cmd/internal/database"
	"log"
	"testing"
	"time"
)

var db *sql.DB
var m *SnippetModel
var s Snippet

func init() {
	db = database.New()
	m = &SnippetModel{DB: db}
	s = Snippet{
		ID:      0,
		Title:   "Testowy snippet",
		Content: "Testowy content",
		Created: time.Now(),
		Expires: "2025-05-05",
	}

}

// func TestTinsert(t *testing.T) {

// 	id, err := m.Insert(&s)
// 	if err != nil {
// 		t.Errorf("Error inserting snippet: %v", err)
// 	}

// 	incomingId <- *id

// }

// func TestGet(t *testing.T) {

// 	s, err := m.Get(<-incomingId)
// 	if err != nil {
// 		t.Errorf("Error getting snippet: %v", err)
// 	}

// 	log.Printf("%#v", s)

// }

func TestSnippetOperations(t *testing.T) {
	var snippetID *int

	t.Run("Insert", func(t *testing.T) {
		s := Snippet{
			ID:      0,
			Title:   "Testowy snippet",
			Content: "Testowy content",
			Created: time.Now(),
			Expires: "2027-05-05",
		}

		id, err := m.Insert(&s)
		if err != nil {
			t.Fatalf("Error inserting: %v", err)
		}
		snippetID = id
	})

	t.Run("Get", func(t *testing.T) {
		if snippetID == nil {
			t.Skip("Skipping Get test, Insert failed")
		}

		var jeden = *snippetID

		s, err := m.Get(jeden)
		if err != nil {
			t.Fatalf("Error getting: %v", err)
		}
		log.Printf("%#v", s)
	})
}
