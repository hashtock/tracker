package core

import (
	"net/http"
)

type Who interface {
	Who(req *http.Request) (*User, error)
}
