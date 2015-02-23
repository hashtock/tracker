package core

import (
    "time"
)

type CountTracker interface {
    AddTagCounts(tagCounts []TagCount) error
}

type CountReader interface {
    Tags() ([]Tag, error)
}

type CountReaderLast interface {
    CountsLast(duration time.Duration) ([]TagCount, error)
    TrendsLast(duration time.Duration) ([]TagCountTrend, error)
}

type CountReaderSince interface {
    CountsSince(since time.Time) ([]TagCount, error)
    TrendsSince(since time.Time) ([]TagCountTrend, error)
}

type CountWritter interface {
    AddTag(tag string) error
}

type CountDestroyer interface {
    RemoveCounts() error
    RemoveAll() error
}

type Counter interface {
    CountDestroyer
    CountReader
    CountReaderLast
    CountReaderSince
    CountTracker
    CountWritter
}

type CountReaderWritter interface {
    CountReader
    CountReaderLast
    CountReaderSince
    CountWritter
}
