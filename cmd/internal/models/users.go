package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id             int
	Name           string
	Email          string
	Password       string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(name, email string, Password string) (*int, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(Password), 12)
	if err != nil {
		return nil, err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created) 
		VALUES ($1, $2, $3, NOW()) RETURNING id;`
	var id int
	err = m.DB.QueryRow(stmt, name, email, hashedPassword).Scan(&id)
	if err != nil {
		//
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			// 23505 = unique_violation
			if pqErr.Code == "23505" {
				return nil, ErrDuplicateEmail // your custom error
			}
		}
		fmt.Printf("Error inserting user: %v\n", err)
		return nil, fmt.Errorf("%#v", err)

	}

	return &id, nil

}

func (m *UserModel) Authenticate(email string, hashedPassword []byte) (*User, error) {
	stmt := `SELECT id, name, email, hashed_password, created FROM users WHERE email = $1;`
	var user User
	err := m.DB.QueryRow(stmt, email).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Created)
	if err != nil {
		return nil, err
	}
	return &user, nil

}

func (m *UserModel) Exists(id int) (bool, error) {
	stmt := `SELECT id, name, email, hashed_password, created FROM users WHERE email = $1;`
	var user User
	err := m.DB.QueryRow(stmt, id).Scan(&user.Id, &user.Name, &user.Email, &user.HashedPassword, &user.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
