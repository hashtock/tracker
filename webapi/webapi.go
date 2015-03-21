package webapi

import (
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"

	"github.com/hashtock/tracker/core"
)

func RunWebApi(counter core.CountReaderWritter) {
	hmacAuth := newVanGoh()

	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if false {
			hmacAuth.ChainedHandler(res, req, nil)
		}
	})

	cs := counterService{counter}

	m.Group("/api", func(r martini.Router) {
		r.Group("/tag", func(sr martini.Router) {
			sr.Get("/", cs.allTags)
			sr.Put("/:name/", cs.addTag)
		})
		r.Get("/counts", cs.counts)
		r.Get("/trends", cs.trends)
	})

	m.Run()
}
