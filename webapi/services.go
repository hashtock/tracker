package webapi

import (
    "errors"
    "net/http"
    "net/url"
    "time"

    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"

    "github.com/hashtock/tracker/core"
)

type counterService struct {
    counter core.CountReaderWritter
}

// Using since, until and duration from query figures out actual since and util
// Valid queries:
// since, until    -> since,          until
// duration        -> now - duration, now
// since           -> since,          now
// until           -> 0,              until
// nothing         -> now - 24h,      now
// dates, duration -> error
func parseQuery(query url.Values) (since, until time.Time, err error) {
    var duration time.Duration

    since_str := query.Get("since")
    until_str := query.Get("until")
    duration_str := query.Get("duration")

    if (since_str != "" || until_str != "") && duration_str != "" {
        err = errors.New("Dates and duration specified")
        return
    }

    if duration_str != "" {
        if duration, err = time.ParseDuration(duration_str); err != nil {
            return
        }
    }

    if since_str != "" {
        if since, err = time.Parse(time.RFC3339, since_str); err != nil {
            return
        }
    }

    if until_str != "" {
        if until, err = time.Parse(time.RFC3339, until_str); err != nil {
            return
        }
    }

    if since.IsZero() && until.IsZero() && duration == 0 {
        duration = time.Hour * 24
    }

    if until.IsZero() {
        until = time.Now()
    }

    if duration != 0 {
        until = time.Now()
        since = until.Add(-duration)
    }

    return
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
