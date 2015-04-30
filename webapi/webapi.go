package webapi

import (
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"

	"github.com/hashtock/tracker/core"
)

func RunWebAPI(counter core.CountReaderWritter, serializer Serializer) {
	hmacAuth := newVanGoh()

	n := negroni.New(
		negroni.NewRecovery(),
		negroni.NewLogger(),
		negroni.HandlerFunc(hmacAuth.ChainedHandler),
	)

	cs := counterService{counter, serializer}

	r := mux.NewRouter()
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/tag/", cs.allTags).Methods("GET")
	api.HandleFunc("/tag/{name}/", cs.addTag).Methods("PUT")
	api.HandleFunc("/counts/", cs.counts).Methods("GET")
	api.HandleFunc("/trends/", cs.trends).Methods("GET")
	api.HandleFunc("/trends/{name}/", cs.tagTrends).Methods("GET")

	n.UseHandler(r)
	n.Run(":3001")
}
