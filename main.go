package main

import (
    "log"
    "net/url"
    "time"

    "code.google.com/p/gcfg"
    "github.com/ChimeraCoder/anaconda"
)

type Config struct {
    Auth struct {
        ConsumerKey       string
        SecretKey         string
        AccessToken       string
        AccessTokenSecret string
    }
    General struct {
        Timeout time.Duration
    }
    Tags struct {
        Tags []string
    }
}

var cfg Config

func init() {
    err := gcfg.ReadFileInto(&cfg, "config.ini")
    if err != nil {
        log.Panicln("Config error:", err.Error())
    }

    example_config := `[auth]
        ConsumerKey = "123"
        SecretKey   = "456"
        AccessToken = "679"
        AccessTokenSecret = "001"
    `

    if cfg.Auth.ConsumerKey == "" || cfg.Auth.SecretKey == "" || cfg.Auth.AccessToken == "" || cfg.Auth.AccessTokenSecret == "" {
        log.Panicln("Twitter authentication missing!\nExpect:", example_config)
    }
}

func getApi() *anaconda.TwitterApi {
    anaconda.SetConsumerKey(cfg.Auth.ConsumerKey)
    anaconda.SetConsumerSecret(cfg.Auth.SecretKey)
    return anaconda.NewTwitterApi(cfg.Auth.AccessToken, cfg.Auth.AccessTokenSecret)
}

func getTags() []string {
    return cfg.Tags.Tags
}

func tagsToTrack() (tags []string) {
    for _, tag := range getTags() {
        tags = append(tags, "#"+tag)
    }
    return tags
}

func getTagStream(api *anaconda.TwitterApi, tags []string) anaconda.Stream {
    values := make(url.Values)
    values["track"] = tags

    stream, err := api.PublicStreamFilter(values)
    if err != nil {
        log.Panicf("Could not get stream: %v", err.Error())
    }

    return stream
}

func main() {
    api := getApi()
    tags := tagsToTrack()
    stream := getTagStream(api, tags)

    go func() {
        time.Sleep(time.Second * cfg.General.Timeout)
        stream.Close()
    }()

    counts := make(map[string]int)
    for _, tag := range getTags() {
        counts[tag] = 0
    }

    log.Println("Start")

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

    log.Println("OK")
    log.Println(counts)
}
