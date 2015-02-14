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

func CliApp() *cli.App {
    app := cli.NewApp()
    app.Name = "tracker"
    app.Usage = "Twitter hashtag count tracking"
    app.Author = "Karol Dulęba"
    app.Email = "mr.fuxi@gmail.com"
    app.Version = "0.1"
    app.Flags = []cli.Flag{
        cli.BoolFlag{Name: "verbose", Usage: "be more verbose"},
    }
    app.Action = cmdWebApi

    app.Commands = []cli.Command{
        {
            Name:      "web",
            ShortName: "w",
            Usage:     "run web server",
            Action:    cmdWebApi,
        },
        {
            Name:      "listen",
            ShortName: "l",
            Usage:     "start counting tweets with requested hashtags",
            Action:    cmdListen,
        },
        {
            Name:      "tags",
            ShortName: "t",
            Usage:     "list current tags",
            Action:    cmdListTags,
        },
        {
            Name:      "counts",
            ShortName: "c",
            Usage:     "list counts for last hour - sum",
            Action:    cmdListTagCounts,
        },
        {
            Name:      "counts detailed",
            ShortName: "cd",
            Usage:     "list counts for last hour - all data points",
            Action:    cmdListTagCountsDetails,
        },
        {
            Name:      "add",
            ShortName: "a",
            Usage:     "list current tags",
            Action:    cmdAddTags,
        },
        {
            Name:   "clear_all",
            Usage:  "removes all data",
            Action: cmdClearAll,
        },
        {
            Name:   "clear_counts",
            Usage:  "removes all counts",
            Action: cmdClearCounts,
        },
    }
    return app
}

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

func cmdListTags(ctx *cli.Context) {
    tags := storage.GetTagsToTrack()

    if len(tags) == 0 {
        fmt.Println("No tags in the system")
        return
    }

    fmt.Println("Tags in system:")
    for _, tag := range tags {
        fmt.Println(tag.Name)
    }
}

func cmdListTagCounts(ctx *cli.Context) {
    tagCounts := storage.GetTagCountForLast(time.Hour * 1)

    if len(tagCounts) == 0 {
        fmt.Println("No tag counts in the system")
        return
    }

    fmt.Println("Counts for last hour:")
    for _, tag := range tagCounts {
        fmt.Printf("%v - %v\n", tag.Name, tag.Count)
    }
}

func cmdListTagCountsDetails(ctx *cli.Context) {
    tagCountTrend := storage.GetTagDetailedCountForLast(time.Hour * 1)

    if len(tagCountTrend) == 0 {
        fmt.Println("No tag counts in the system")
        return
    }

    fmt.Println("Detailed counts for last hour:")
    for _, tag := range tagCountTrend {
        fmt.Println(tag.Name)
        for _, datapoint := range tag.Counts {
            fmt.Printf("\t%v - %v\n", datapoint.Date, datapoint.Count)
        }
    }
}

func cmdAddTags(ctx *cli.Context) {
    tags := append([]string{ctx.Args().First()}, ctx.Args().Tail()...)
    fmt.Println("Tags to add:", tags)
    if err := storage.AddTagsToTrack(tags); err != nil {
        fmt.Println("Could not store new tags.", err.Error())
    }
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

func cmdWebApi(ctx *cli.Context) {
    cfg := conf.GetConfig()
    cfg.General.Timeout = "0" // No timeout
    go cmdListen(ctx)

    webapi.RunWebApi()
}
