package core

import (
    "time"
)

type CountTracker interface {
    AddTagCounts(tagCounts []TagCount) error
}

type CountReader interface {
    Tags() ([]Tag, error)
    CountsLast(duration time.Duration) ([]TagCount, error)
    TrendsLast(duration time.Duration) ([]TagCountTrend, error)
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
    CountTracker
    CountWritter
}

type CountReaderWritter interface {
    CountReader
    CountWritter
}
