package gateway

import (
	"errors"
	"os"
)

type UserCredentials struct {
	Username string
	Password string
}

func GetEnv() (*UserCredentials, error) {
	username := os.Getenv("username")
	password := os.Getenv("password")

	if username == "" || password == "" {
		return nil, errors.New("username or password not set")
	}

	return &UserCredentials{
		Username: username,
		Password: password,
	}, nil
}
