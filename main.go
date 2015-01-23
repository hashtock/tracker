package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "code.google.com/p/gcfg"
    "github.com/codegangsta/cli"
)

type auth struct {
    ConsumerKey       string
    SecretKey         string
    AccessToken       string
    AccessTokenSecret string
}

type config struct {
    Auth    auth
    General struct {
        Timeout time.Duration
        DB      string
    }
}

var cfg config

func init() {
    err := gcfg.ReadFileInto(&cfg, "config.ini")
    if err != nil {
        log.Fatalln("Config error:", err.Error())
    }

    example_config := `[auth]
        ConsumerKey = "123"
        SecretKey   = "456"
        AccessToken = "679"
        AccessTokenSecret = "001"
    `

    if cfg.Auth.ConsumerKey == "" || cfg.Auth.SecretKey == "" || cfg.Auth.AccessToken == "" || cfg.Auth.AccessTokenSecret == "" {
        log.Fatalln("Twitter authentication missing!\nExpect:", example_config)
    }

    if err := startSession(cfg.General.DB); err != nil {
        log.Fatalln("Could not connect to DB.", err.Error())
    }
}

func cmdListen(ctx *cli.Context) {
    now := time.Now()

    tags := GetTagsToTrack()
    tagNames := make([]string, len(tags))
    for i, tag := range tags {
        tagNames[i] = tag.Name
    }
    countMap := Listen(tagNames, cfg.General.Timeout, cfg.Auth)
    tc := make([]TagCount, 0, len(countMap))

    for tagName, count := range countMap {
        tc = append(tc, TagCount{
            Name:  tagName,
            Count: count,
            Date:  now,
        })
    }

    if err := AddTagCounts(tc); err != nil {
        fmt.Println("Could not store tag counts.", err.Error())
    }
    fmt.Println("Done")
}

func cmdListTags(ctx *cli.Context) {
    tags := GetTagsToTrack()

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
    tagCounts := GetTagCountFor(time.Hour * 1)

    if len(tagCounts) == 0 {
        fmt.Println("No tag counts in the system")
        return
    }

    fmt.Println("Counts for last hour:")
    for _, tag := range tagCounts {
        fmt.Printf("%v - %v\n", tag.Name, tag.Count)
    }
}

func cmdAddTags(ctx *cli.Context) {
    tags := append([]string{ctx.Args().First()}, ctx.Args().Tail()...)
    fmt.Println("Tags to add:", tags)
    if err := AddTagsToTrack(tags); err != nil {
        fmt.Println("Could not store new tags.", err.Error())
    }
}

func cmdClearAll(ctx *cli.Context) {
    if err := DropAll(); err != nil {
        fmt.Println("Could not drop all data.", err.Error())
    }
    fmt.Println("Done. Nothing left...")
}

func cmdClearCounts(ctx *cli.Context) {
    if err := DropCollection(TAG_COUNT_COLLECTION); err != nil {
        fmt.Println("Could not drop counts.", err.Error())
    }
    fmt.Println("Done. Counts cleared tags are still in place")
}

func main() {
    app := cli.NewApp()
    app.Name = "tracker"
    app.Usage = "Twitter hashtag count tracking"
    app.Author = "Karol Dulęba"
    app.Email = "mr.fuxi@gmail.com"
    app.Version = "0.1"
    app.Action = func(c *cli.Context) {
        cli.ShowAppHelp(c)
    }

    app.Commands = []cli.Command{
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
            Usage:     "list counts for last hour",
            Action:    cmdListTagCounts,
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

    app.Run(os.Args)
}
