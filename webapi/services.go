package webapi

import (
    "net/http"
    "time"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/tracker/storage"
)

func allTags(r render.Render) {
    tags := storage.GetTagsToTrack()
    r.JSON(http.StatusOK, tags)
}

func addTag(res http.ResponseWriter, params martini.Params) {
    storage.AddTagsToTrack([]string{params["name"]})
    res.WriteHeader(http.StatusCreated)
}

func countLastDay(r render.Render) {
    today := time.Now().Truncate(time.Hour * 24)
    yesterday := today.Add(-time.Hour * 24)

    counts := storage.GetTagCount(yesterday, today)
    r.JSON(http.StatusOK, counts)
}

func countForDuration(params martini.Params, r render.Render) {
    duration_str, ok := params["duration"]
    if !ok {
        duration_str = "24h"
    }

    duration, err := time.ParseDuration(duration_str)
    if err != nil {
        r.Error(http.StatusBadRequest)
        return
    }

    counts := storage.GetTagCountForLast(duration)
    r.JSON(http.StatusOK, counts)
}
