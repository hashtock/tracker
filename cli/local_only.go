package cli

import (
    "fmt"
    "log"
    "time"

    "github.com/codegangsta/cli"

    "github.com/hashtock/tracker/conf"
    "github.com/hashtock/tracker/core"
    "github.com/hashtock/tracker/listener"
    "github.com/hashtock/tracker/webapi"
)

func cmdListen(ctx *cli.Context) {
    cfg := conf.GetConfig()
    counter := getCounter(ctx)

    tags, err := counter.Tags()
    if err != nil {
        log.Fatalln(err)
    }

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
        tc := make([]core.TagCount, 0, len(countMap))

        for tagName, count := range countMap {
            tc = append(tc, core.TagCount{
                Name:  tagName,
                Count: count,
                Date:  now,
            })
        }

        if err := counter.AddTagCounts(tc); err != nil {
            fmt.Println("Could not store tag counts.", err.Error())
        }
    }
}

func cmdWebApi(ctx *cli.Context) {
    counter := getCounterRW(ctx)
    cfg := conf.GetConfig()
    cfg.General.Timeout = "0" // No timeout
    go cmdListen(ctx)

    webapi.RunWebApi(counter)
}

func cmdClearAll(ctx *cli.Context) {
    counter := getCounter(ctx)

    if err := counter.RemoveAll(); err != nil {
        log.Fatalln("Could not remove the data: ", err)
    }
    fmt.Println("Done. Nothing left...")
}

func cmdClearCounts(ctx *cli.Context) {
    counter := getCounter(ctx)

    if err := counter.RemoveCounts(); err != nil {
        log.Println("Could not remove counts.", err)
    }
    fmt.Println("Done. Counts cleared tags are still in place")
}
