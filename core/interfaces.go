package core

import (
	"time"
)

type DataNotificator interface {
	DataAvailable(since, until time.Time)
}

type MessagePublisher interface {
	Publish(subject string, v interface{}) error
}
