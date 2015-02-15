package cli

import (
    "fmt"
    "log"

    "github.com/codegangsta/cli"

    "github.com/hashtock/tracker/core"
    "github.com/hashtock/tracker/storage"
)

func cmdListTags(ctx *cli.Context) {
    var tags []core.Tag
    var err error

    if tracker := getRemoteClient(ctx); tracker != nil {
        tags, err = tracker.GetTagList()
        if err != nil {
            log.Fatalln(err)
        }
    } else {
        tags = storage.GetTagsToTrack()
    }

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
    var tagCounts []core.TagCount
    var err error

    timeSpan := getDuration(ctx)

    if tracker := getRemoteClient(ctx); tracker != nil {
        tagCounts, err = tracker.GetTagCounts(timeSpan)
        if err != nil {
            log.Fatalln(err)
        }
    } else {
        tagCounts = storage.GetTagCountForLast(timeSpan)
    }

    if len(tagCounts) == 0 {
        fmt.Println("No tag counts in the system")
        return
    }

    fmt.Printf("Counts for last %s:\n", timeSpan)
    for _, tag := range tagCounts {
        fmt.Printf("%v - %v\n", tag.Name, tag.Count)
    }
}

func cmdListTagCountsDetails(ctx *cli.Context) {
    var tagCountTrend []core.TagCountTrend
    var err error

    timeSpan := getDuration(ctx)

    if tracker := getRemoteClient(ctx); tracker != nil {
        tagCountTrend, err = tracker.GetTagTrends(timeSpan)
        if err != nil {
            log.Fatalln(err)
        }
    } else {
        tagCountTrend = storage.GetTagDetailedCountForLast(timeSpan)
    }

    if len(tagCountTrend) == 0 {
        fmt.Println("No tag counts in the system")
        return
    }

    fmt.Printf("Detailed counts for last %s:\n", timeSpan)
    for _, tag := range tagCountTrend {
        fmt.Println(tag.Name)
        for _, datapoint := range tag.Counts {
            fmt.Printf("\t%v - %v\n", datapoint.Date, datapoint.Count)
        }
    }
}

func cmdAddTags(ctx *cli.Context) {
    var err error

    tags := append([]string{ctx.Args().First()}, ctx.Args().Tail()...)

    fmt.Println("Tags to add:", tags)
    if tracker := getRemoteClient(ctx); tracker != nil {
        for _, tag := range tags {
            if err = tracker.AddTag(tag); err != nil {
                break
            }
        }
    } else {
        err = storage.AddTagsToTrack(tags)
    }

    if err != nil {
        fmt.Println("Could not store new tags.", err.Error())
    }
}
