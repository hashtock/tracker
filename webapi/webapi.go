package webapi

import (
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
)

func RunWebApi() {
    hmacAuth := newVanGoh()

    m := martini.Classic()
    m.Use(render.Renderer())
    m.Use(func(res http.ResponseWriter, req *http.Request) {
        hmacAuth.ChainedHandler(res, req, nil)
    })

    m.Group("/api/tag", func(r martini.Router) {
        r.Get("/", allTags)
        r.Put("/:name/", addTag)
    })
    m.Group("/api/counts", func(r martini.Router) {
        r.Get("/", countLastDay)
        r.Get("/:duration/", countForDuration)
    })
    m.Run()
}
