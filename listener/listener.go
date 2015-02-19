package listener

import (
    "net/url"
    "time"

    "github.com/ChimeraCoder/anaconda"

    "github.com/hashtock/tracker/conf"
)

type Listener interface {
    Listen() chan map[string]int
}

type twitterListener struct {
    timeout time.Duration
    update  time.Duration
    auth    conf.Auth
    tags    []string

    api         *anaconda.TwitterApi
    stream      anaconda.Stream
    counter     *tagCounter
    dataChannel chan map[string]int

    timeOutCh chan bool
}

func NewTwitterListener(tags []string, timeout time.Duration, update time.Duration, auth conf.Auth) Listener {
    listener := &twitterListener{
        timeout: timeout,
        update:  update,
        auth:    auth,
        tags:    tags,
    }

    listener.counter = newCounter()
    listener.dataChannel = make(chan map[string]int, 0)

    listener.timeOutCh = make(chan bool)

    return listener
}

func (t *twitterListener) connectToApi() {
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
        t.connectToApi()
    }

    t.stream = t.api.PublicStreamFilter(values)
}

func (t *twitterListener) keepUpdatingClientWithData() {
    timer := time.NewTicker(t.update)
    for {
        select {
        case <-timer.C:
            t.dataChannel <- t.counter.getDataAndClear()
        case <-t.timeOutCh:
            timer.Stop()
            return
        }
    }
}

func (t *twitterListener) watchForRunningTimout() {
    if t.timeout <= 0 {
        return
    }

    time.Sleep(t.timeout)
    close(t.timeOutCh)
    t.stream.Interrupt()
    t.dataChannel <- t.counter.getDataAndClear()
    close(t.dataChannel)
}

func (t *twitterListener) processTweets() {
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
    go t.keepUpdatingClientWithData()
    go t.watchForRunningTimout()

    t.startTagStream()
    go t.processTweets()

    return t.dataChannel
}
