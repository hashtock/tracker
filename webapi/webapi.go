package webapi

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
)

func RunWebApi() {
    m := martini.Classic()
    m.Use(render.Renderer())

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
