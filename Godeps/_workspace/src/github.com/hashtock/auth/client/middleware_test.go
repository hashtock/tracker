package client_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"

	"github.com/hashtock/auth/client"
	"github.com/hashtock/auth/core"
)

type WhoMock struct {
	Err error
}

var testUser = &core.User{
	Name:   "Bob",
	Email:  "bob@example.com",
	Avatar: "http://example.com/bob.png",
}

func (w WhoMock) Who(req *http.Request) (user *core.User, err error) {
	if w.Err == nil {
		user = testUser
	}
	err = w.Err
	return
}

func TestNegroniMiddleware(t *testing.T) {
	middleware := client.NewAuthMiddleware(WhoMock{})
	assert.Implements(t, (*negroni.Handler)(nil), middleware)
}

func TestWhoSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/some/url", nil)

	middleware := client.NewAuthMiddleware(WhoMock{Err: nil})
	n := negroni.New(middleware)
	n.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		user := context.Get(r, client.UserContextKey)
		if assert.NotNil(t, user) {
			assert.EqualValues(t, testUser, user)
		}

		rw.WriteHeader(http.StatusTeapot)
	})
	n.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusTeapot, w.Code)
}

func TestWhoNotLoggedIn(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/some/url", nil)

	middleware := client.NewAuthMiddleware(WhoMock{Err: core.ErrUserNotLoggedIn})
	n := negroni.New(middleware)
	n.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "Handler should not be called")
	})
	n.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusUnauthorized, w.Code)
}

func TestWhoGenericError(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/some/url", nil)

	middleware := client.NewAuthMiddleware(WhoMock{Err: errors.New("Error!")})
	n := negroni.New(middleware)
	n.UseHandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		assert.Fail(t, "Handler should not be called")
	})
	n.ServeHTTP(w, req)

	assert.EqualValues(t, http.StatusInternalServerError, w.Code)
}
