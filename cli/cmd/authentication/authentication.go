package authentication

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	id uuid.UUID
	username string
	password string
}

var (
	currentUser *User
)

func login(username, password string) error {
	if username != "" && password != "" {
		currentUser = &User{
			id: uuid.New(),
			username: username,
			password: password,
		}
	} else {
		return errors.New("provided username or password is invalid")
	}

	return nil
}

func logout() error {
	currentUser = nil
	return nil
}

func register(username, password string) error {
	// For now it is the same as login
	return login(username, password)
}