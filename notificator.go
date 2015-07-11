package main

import (
	"log"
	"time"

	"github.com/hashtock/tracker/core"
)

type MsgNotificator struct {
	counter      core.CountReader
	msgPublisher core.MessagePublisher
}

func NewMsgNoticiator(counter core.CountReader, msgPublisher core.MessagePublisher) *MsgNotificator {
	return &MsgNotificator{
		counter:      counter,
		msgPublisher: msgPublisher,
	}
}

func (m *MsgNotificator) DataAvailable(since, until time.Time) {
	tags, err := m.counter.Tags()
	if err != nil {
		log.Println("Error getting list of tags:", err)
		return
	}
	counts, err := m.counter.Counts(since, until)
	if err != nil {
		log.Println("Error getting recent counts:", err)
		return
	}

	data := make(map[string]int, len(tags))
	for _, tag := range tags {
		data[tag.Name] = 0
	}
	for _, tagCount := range counts {
		data[tagCount.Name] += tagCount.Count
	}

	m.msgPublisher.Publish("tags.counts", data)
}
