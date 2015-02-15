package cli

import (
    "fmt"
    "time"

    "github.com/codegangsta/cli"

    "github.com/hashtock/tracker/conf"
    "github.com/hashtock/tracker/listener"
    "github.com/hashtock/tracker/storage"
    "github.com/hashtock/tracker/webapi"
)

func cmdListen(ctx *cli.Context) {
    cfg := conf.GetConfig()

    tags := storage.GetTagsToTrack()
    tagNames := make([]string, len(tags))
    for i, tag := range tags {
        tagNames[i] = tag.Name
    }

    countCh := listener.Listen(tagNames, cfg.General.TimeoutD(), cfg.General.UpdateTimeD(), cfg.Auth)

    for countMap := range countCh {
        now := time.Now().Truncate(cfg.General.SampingTimeD())
        if ctx.GlobalBool("verbose") {
            fmt.Printf("Time: %v\tData: %v\n", now, countMap)
        }
        tc := make([]storage.TagCount, 0, len(countMap))

        for tagName, count := range countMap {
            tc = append(tc, storage.TagCount{
                Name:  tagName,
                Count: count,
                Date:  now,
            })
        }

        if err := storage.AddTagCounts(tc); err != nil {
            fmt.Println("Could not store tag counts.", err.Error())
        }
    }
}

func cmdWebApi(ctx *cli.Context) {
    cfg := conf.GetConfig()
    cfg.General.Timeout = "0" // No timeout
    go cmdListen(ctx)

    webapi.RunWebApi()
}

func cmdClearAll(ctx *cli.Context) {
    if err := storage.DropAll(); err != nil {
        fmt.Println("Could not drop all data.", err.Error())
    }
    fmt.Println("Done. Nothing left...")
}

func cmdClearCounts(ctx *cli.Context) {
    if err := storage.DropCollection(storage.TAG_COUNT_COLLECTION); err != nil {
        fmt.Println("Could not drop counts.", err.Error())
    }
    fmt.Println("Done. Counts cleared tags are still in place")
}
