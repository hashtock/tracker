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

    m.Group("/api/tag", func(r martini.Router) {
        r.Get("/", cs.allTags)
        r.Put("/:name/", cs.addTag)
    })
    m.Group("/api/counts", func(r martini.Router) {
        r.Get("/", cs.countForDuration)
        r.Get("/:duration/", cs.countForDuration)
        r.Get("/since/:since/", cs.countSince)
    })
    m.Group("/api/trends", func(r martini.Router) {
        r.Get("/", cs.countDetailsForDuration)
        r.Get("/:duration/", cs.countDetailsForDuration)
        r.Get("/since/:since/", cs.countDetailsSince)
    })

    m.Run()
}
