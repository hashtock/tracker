package main

import (
    "log"
    "net/url"
    "time"

    "github.com/ChimeraCoder/anaconda"
)

func getApi(twAuth auth) *anaconda.TwitterApi {
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

func Listen(tags []string, timeout time.Duration, twAuth auth) (counts map[string]int) {
    api := getApi(twAuth)
    hashedTags := tagsToTrack(tags)
    stream := getTagStream(api, hashedTags)

    go func() {
        time.Sleep(time.Second * timeout)
        close(stream.C)
    }()

    counts = make(map[string]int)
    for _, tag := range tags {
        counts[tag] = 0
    }

    for msg := range stream.C {
        tweet, ok := msg.(anaconda.Tweet)
        if !ok {
            continue
        }

        tags := tweet.Entities.Hashtags
        for _, tag := range tags {
            if _, ok := counts[tag.Text]; ok {
                counts[tag.Text]++
            }
        }
    }

    return
}
