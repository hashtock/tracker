package webapi

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/context"
	authCore "github.com/hashtock/auth/core"

	"github.com/hashtock/tracker/core"
)

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

	sinceStr := query.Get("since")
	untilStr := query.Get("until")
	durationStr := query.Get("duration")

	if (sinceStr != "" || untilStr != "") && durationStr != "" {
		err = errors.New("Dates and duration specified")
		return
	}

	if durationStr != "" {
		if duration, err = time.ParseDuration(durationStr); err != nil {
			return
		}
	}

	if sinceStr != "" {
		if since, err = time.Parse(time.RFC3339, sinceStr); err != nil {
			return
		}
	}

	if untilStr != "" {
		if until, err = time.Parse(time.RFC3339, untilStr); err != nil {
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

func getSamplingFromQuery(query url.Values) (sampling core.Sampling, err error) {
	samplingStr := query.Get("sampling")
	if samplingStr == "" {
		return core.SamplingRaw, nil
	}

	return core.ParseSampling(samplingStr)
}

func isAdmin(req *http.Request) error {
	obj := context.Get(req, UserContextKey)
	if obj == nil {
		return errors.New("User not logged in")
	}

	user, ok := obj.(authCore.User)
	if !ok || !user.Admin {
		return errors.New("Admin permission needed")
	}

	return nil
}
