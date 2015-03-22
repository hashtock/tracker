package listener

import (
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/ChimeraCoder/anaconda"

	"github.com/hashtock/tracker/conf"
)

type twitterListener struct {
	timeout time.Duration
	update  time.Duration
	auth    conf.Auth
	tags    []string

	api         *anaconda.TwitterApi
	stream      anaconda.Stream
	counter     *tagCounter
	dataChannel chan map[string]int

	stopAll chan struct{}
	allDone sync.WaitGroup
}

func NewTwitterListener(tags []string, timeout time.Duration, update time.Duration, auth conf.Auth) Listener {
	listener := &twitterListener{
		timeout: timeout,
		update:  update,
		auth:    auth,
		tags:    tags,
	}

	listener.counter = newCounter()
	listener.dataChannel = make(chan map[string]int, 100)

	return listener
}

func (t *twitterListener) connectToAPI() {
	anaconda.SetConsumerKey(t.auth.ConsumerKey)
	anaconda.SetConsumerSecret(t.auth.SecretKey)
	t.api = anaconda.NewTwitterApi(t.auth.AccessToken, t.auth.AccessTokenSecret)
}

func (t *twitterListener) startTagStream() {
	hashedTags := make([]string, len(t.tags))
	for i, tag := range t.tags {
		hashedTags[i] = "#" + tag
	}

	values := make(url.Values)
	values["track"] = hashedTags

	if t.api == nil {
		t.connectToAPI()
	}

	t.stream = t.api.PublicStreamFilter(values)
}

func (t *twitterListener) keepUpdatingClientWithData() {
	ticker := time.NewTicker(t.update)
	t.allDone.Add(1)
	defer t.allDone.Done()

	for {
		select {
		case <-ticker.C:
			t.dataChannel <- t.counter.getDataAndClear()
		case <-t.stopAll:
			ticker.Stop()
			return
		}
	}
}

func (t *twitterListener) watchForRunningTimout() {
	if t.timeout <= 0 {
		return
	}

	t.allDone.Add(1)
	defer t.allDone.Done()

	timeout := time.NewTimer(t.timeout)
	for {
		select {
		case <-timeout.C:
			t.Stop()
		case <-t.stopAll:
			timeout.Stop()
			return
		}
	}
}

func (t *twitterListener) processTweets() {
	t.allDone.Add(1)
	defer t.allDone.Done()

	tagsMap := make(map[string]bool)
	for _, tag := range t.tags {
		tagsMap[tag] = true
	}

	for msg := range t.stream.C {
		tweet, ok := msg.(anaconda.Tweet)
		if !ok {
			continue
		}

		tags := tweet.Entities.Hashtags
		for _, tag := range tags {
			if _, ok := tagsMap[tag.Text]; ok {
				t.counter.incCount(tag.Text, 1)
			}
		}
	}
}

func (t *twitterListener) Listen() chan map[string]int {
	if t.stopAll != nil {
		return t.dataChannel
	}

	t.stopAll = make(chan struct{})

	go t.keepUpdatingClientWithData()
	go t.watchForRunningTimout()

	t.startTagStream()
	go t.processTweets()

	return t.dataChannel
}

func (t *twitterListener) Stop() {
	if t.stopAll == nil {
		return
	}

	close(t.stopAll)

	t.stream.Interrupt()
	close(t.stream.C)

	t.dataChannel <- t.counter.getDataAndClear()
	close(t.dataChannel)

	t.allDone.Wait()
	t.stopAll = nil
}

func (t *twitterListener) SetTags(tags []string) {
	t.tags = tags

	// Stop here if stream is not running
	if t.stream.Quit == nil {
		log.Println("Stream is not running yet")
		return
	}

	// Close old stream
	log.Println("Closing old stream")
	t.stream.Interrupt()
	close(t.stream.C)

	// Start new stream and start processing new tags
	log.Println("Starting new stream")
	t.startTagStream()
	go t.processTweets()
}

func (t *twitterListener) Tags() []string {
	return t.tags
}
