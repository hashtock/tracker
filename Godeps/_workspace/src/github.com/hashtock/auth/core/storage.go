package core

import (
	"errors"
)

var ErrSessionNotFound error = errors.New("Session not found")

type UserSessioner interface {
	GetUserBySession(sessionId string) (*User, error)
	AddUserToSession(sessionId string, user *User) error
	DeleteSession(sessionId string) error
}

type Administrator interface {
	MakeUserAnAdmin(email string) error
}
