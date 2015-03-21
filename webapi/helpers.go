package webapi

import (
	"errors"
	"net/url"
	"time"
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
