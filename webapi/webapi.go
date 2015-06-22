package webapi

import (
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/pat"
	authClient "github.com/hashtock/auth/client"
	authCore "github.com/hashtock/auth/core"

	"github.com/hashtock/tracker/core"
)

type Options struct {
	Counter    core.CountReaderWritter
	Serializer Serializer
	WhoClient  authCore.Who
}

func Handlers(options Options) http.Handler {
	hmacAuth := newVanGoh()

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		authClient.NewAuthMiddleware(options.WhoClient),
		// negroni.HandlerFunc(hmacAuth.ChainedHandler),
	)

	if false {
		n.UseFunc(hmacAuth.ChainedHandler)
	}

	cs := counterService{options.Counter, options.Serializer}

	m := pat.New()
	m.Get("/tag/", cs.allTags)
	m.Put("/tag/{name}/", cs.addTag)
	m.Get("/counts/", cs.counts)
	m.Get("/trends/", cs.trends)
	m.Get("/trends/{name}/", cs.tagTrends)

	n.UseHandler(m)

	return n
}
