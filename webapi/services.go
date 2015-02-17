package webapi

import (
    "net/http"
    "time"

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

func (c *counterService) countForDuration(params martini.Params, r render.Render) {
    duration_str, ok := params["duration"]
    if !ok {
        duration_str = "24h"
    }

    duration, err := time.ParseDuration(duration_str)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    counts, err := c.counter.CountsLast(duration)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    r.JSON(http.StatusOK, counts)
}

func (c *counterService) countDetailsForDuration(params martini.Params, r render.Render) {
    duration_str, ok := params["duration"]
    if !ok {
        duration_str = "24h"
    }

    duration, err := time.ParseDuration(duration_str)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    counts, err := c.counter.TrendsLast(duration)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    r.JSON(http.StatusOK, counts)
}
