package client

import (
	"log"
	"net/http"

	"github.com/gorilla/context"

	"github.com/hashtock/auth/core"
)

const UserContextKey = "user"

type whoMiddleware struct {
	client core.Who
}

func NewAuthMiddleware(who core.Who) *whoMiddleware {
	return &whoMiddleware{
		client: who,
	}
}

func (w whoMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	user, err := w.client.Who(req)
	if err == core.ErrUserNotLoggedIn {
		rw.WriteHeader(http.StatusUnauthorized)
	} else if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		context.Set(req, UserContextKey, user)
		next(rw, req)
		// Gorilla Mux/Pat routers will remove all keys automatically
		context.Delete(req, UserContextKey)
	}
}
