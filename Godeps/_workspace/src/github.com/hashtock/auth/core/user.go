package core

import (
	"errors"
)

var (
	ErrUserNotLoggedIn = errors.New("User not logged in")
)

type User struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Admin  bool   `json:"admin,omitempty"`
}
