package core

import (
	"time"
)

type CountTracker interface {
	AddTagCounts(tagCounts []TagCount) error
}

type CountReader interface {
	Tags() ([]Tag, error)

	Counts(since, until time.Time) ([]TagCount, error)

	Trends(since, until time.Time) ([]TagCountTrend, error)
	TagTrends(tag string, since, until time.Time, sampling Sampling) (TagCountTrend, error)
}

type CountWritter interface {
	AddTag(tag string) error
}

type CountDestroyer interface {
	RemoveCounts() error
	RemoveAll() error
}

type CountReaderWritter interface {
	CountReader
	CountWritter
}

type Counter interface {
	CountDestroyer
	CountReaderWritter
	CountTracker
}

func TagNames(storage CountReader) (tagNames []string, err error) {
	tags, err := storage.Tags()
	if err != nil {
		return
	}

	tagNames = make([]string, len(tags))
	for i, tag := range tags {
		tagNames[i] = tag.Name
	}

	return
}
