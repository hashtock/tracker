package webapi

import (
    "net/http"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/tracker/core"
)

type counterService struct {
    counter core.CountReaderWritter
}

func (c *counterService) allTags(r render.Render) {
    tags, err := c.counter.Tags()
    if err != nil {
        r.Error(http.StatusInternalServerError)
        return
    }

    r.JSON(http.StatusOK, tags)
}

func (c *counterService) addTag(res http.ResponseWriter, params martini.Params) {
    if err := c.counter.AddTag(params["name"]); err != nil {
        res.WriteHeader(http.StatusInternalServerError)
    } else {
        res.WriteHeader(http.StatusCreated)
    }
}

func (c *counterService) counts(req *http.Request, r render.Render) {
    since, until, err := parseQuery(req.URL.Query())
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    counts, err := c.counter.Counts(since, until)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    r.JSON(http.StatusOK, counts)
}

func (c *counterService) trends(req *http.Request, r render.Render) {
    since, until, err := parseQuery(req.URL.Query())
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    trends, err := c.counter.Trends(since, until)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    r.JSON(http.StatusOK, trends)
}
