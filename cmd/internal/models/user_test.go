package models

import (
	"errors"
	"gostats/cmd/internal/assert"
	"gostats/cmd/internal/database"
	"math/rand"
	"strconv"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestInsert(t *testing.T) {
	var um = UserModel{
		DB: database.New(),
	}

	ran := rand.Int()

	passwordStr := "SecretSecret"

	name := "Użytkownik Nr " + strconv.Itoa(ran)
	email := "uzytkownik" + strconv.Itoa(ran) + "@poczta.pl"
	password := passwordStr

	id, err := um.Insert(name, email, password)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := um.User(*id)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, name, expected.Name)
	assert.Equal(t, email, expected.Email)
	println(expected.HashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(expected.HashedPassword), []byte(passwordStr))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			t.Fatal(err)
		} else {
			t.Fatal(err)
		}

	}

}
