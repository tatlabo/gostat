package database

import "testing"

func TestDBConnection(t *testing.T) {
	DB := New()
	err := DB.Ping()
	if err != nil {
		t.Fatal(err)
	}
}
