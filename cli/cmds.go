package cli

import (
    "fmt"
    "log"

    "github.com/codegangsta/cli"
)

func cmdListTags(ctx *cli.Context) {
    counter := getCounterRW(ctx)

    tags, err := counter.Tags()
    if err != nil {
        log.Fatalln(err)
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
    since, until, timeSpan := getTimeRangeFromDuration(ctx)
    counter := getCounterRW(ctx)

    tagCounts, err := counter.Counts(since, until)
    if err != nil {
        log.Fatalln(err)
    }

    if len(tagCounts) == 0 {
        fmt.Println("No tag counts in the system for last", timeSpan)
        return
    }

    fmt.Printf("Counts for last %s:\n", timeSpan)
    for _, tag := range tagCounts {
        fmt.Printf("%v - %v\n", tag.Name, tag.Count)
    }
}

func cmdListTagCountsDetails(ctx *cli.Context) {
    since, until, timeSpan := getTimeRangeFromDuration(ctx)
    counter := getCounterRW(ctx)

    tagCountTrend, err := counter.Trends(since, until)
    if err != nil {
        log.Fatalln(err)
    }

    if len(tagCountTrend) == 0 {
        fmt.Println("No tag counts in the system for last", timeSpan)
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
    counter := getCounterRW(ctx)

    fmt.Println("Tags to add:", tags)
    for _, tag := range tags {
        if err = counter.AddTag(tag); err != nil {
            break
        }
    }

    if err != nil {
        fmt.Println("Could not store new tags.", err.Error())
    }
}
