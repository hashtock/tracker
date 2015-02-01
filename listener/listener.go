package listener

import (
    "log"
    "net/url"
    "time"

    "github.com/ChimeraCoder/anaconda"

    "github.com/hashtock/tracker/conf"
)

func getApi(twAuth conf.Auth) *anaconda.TwitterApi {
    anaconda.SetConsumerKey(twAuth.ConsumerKey)
    anaconda.SetConsumerSecret(twAuth.SecretKey)
    return anaconda.NewTwitterApi(twAuth.AccessToken, twAuth.AccessTokenSecret)
}

func tagsToTrack(tags []string) (hashedTags []string) {
    for _, tag := range tags {
        hashedTags = append(hashedTags, "#"+tag)
    }
    return hashedTags
}

func getTagStream(api *anaconda.TwitterApi, tags []string) anaconda.Stream {
    values := make(url.Values)
    values["track"] = tags

    stream, err := api.PublicStreamFilter(values)
    if err != nil {
        log.Fatalln("Could not get stream: %v", err.Error())
    }

    return stream
}

func Listen(tags []string, timeout time.Duration, update time.Duration, twAuth conf.Auth) (counts chan map[string]int) {
    api := getApi(twAuth)
    hashedTags := tagsToTrack(tags)
    stream := getTagStream(api, hashedTags)
    counter := newCounter()

    counts = make(chan map[string]int, 0)

    stopUpdates := make(chan struct{})
    go func() {
        timer := time.NewTicker(update)
        for {
            select {
            case <-timer.C:
                counts <- counter.getDataAndClear()
            case <-stopUpdates:
                timer.Stop()
                return
            }
        }
    }()

    if timeout > 0 {
        go func() {
            time.Sleep(time.Second * timeout)
            close(stopUpdates)
            stream.Close()
            counts <- counter.getDataAndClear()
        }()
    }

    tagsMap := make(map[string]bool)
    for _, tag := range tags {
        tagsMap[tag] = true
    }

    go func() {
        for msg := range stream.C {
            tweet, ok := msg.(anaconda.Tweet)
            if !ok {
                continue
            }

            tags := tweet.Entities.Hashtags
            for _, tag := range tags {
                if _, ok := tagsMap[tag.Text]; ok {
                    counter.incCount(tag.Text, 1)
                }
            }
        }
    }()

    return counts
}
