package listener

import (
	"log"
	"time"

	"github.com/hashtock/tracker/core"
)

type Listener interface {
	Listen() chan map[string]int
	Stop()

	SetTags(tags []string)
	Tags() []string
}

type Options struct {
	TagUpdateTime time.Duration
	SampingTime   time.Duration
	Verbose       bool
	ExitSignal    chan struct{}
}

func Listen(tagListener Listener, counter core.Counter, options Options) {
	countCh := tagListener.Listen()

	go func() {
		watcher := time.NewTicker(options.TagUpdateTime)
		defer watcher.Stop()

		for {
			select {
			case <-watcher.C:
				newTags, err := core.TagNames(counter)
				if err != nil {
					log.Println("Could not get new tags:", err)
				} else if !compareTags(newTags, tagListener.Tags()) {
					log.Println("Setting new list of tags to track")
					tagListener.SetTags(newTags)
				}
			case <-options.ExitSignal:
				tagListener.Stop()
				return
			}
		}
	}()

	for countMap := range countCh {
		now := time.Now().Truncate(options.SampingTime)
		if options.Verbose {
			log.Printf("Time: %v\tData: %v\n", now, countMap)
		}
		tc := make([]core.TagCount, 0, len(countMap))

		for tagName, count := range countMap {
			tc = append(tc, core.TagCount{
				Name:  tagName,
				Count: count,
				Date:  now,
			})
		}

		if err := counter.AddTagCounts(tc); err != nil {
			log.Println("Could not store tag counts.", err.Error())
		}
	}

	close(options.ExitSignal)
}

func compareTags(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
