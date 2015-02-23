package cli

import (
    "log"
    "time"

    "github.com/codegangsta/cli"

    "github.com/hashtock/tracker/client"
    "github.com/hashtock/tracker/conf"
    "github.com/hashtock/tracker/core"
    "github.com/hashtock/tracker/storage"
)

func getCounterRW(ctx *cli.Context) core.CountReaderWritter {
    var counter core.CountReaderWritter
    var err error

    if remote := ctx.GlobalString("remote"); remote != "" {
        remoteConfig := conf.GetRemoteConfig(remote)
        counter, err = client.NewTracker(remoteConfig)
    } else {
        config := conf.GetConfig()
        counter, err = storage.NewMongoCounter(config.General.DB, config.General.SampingTimeD())
    }

    if err != nil {
        log.Fatalln(err)
    }

    return counter
}

func getCounter(ctx *cli.Context) core.Counter {
    if remote := ctx.GlobalString("remote"); remote != "" {
        log.Fatalln("No remote tracker handles available for that action")
    }

    config := conf.GetConfig()
    counter, err := storage.NewMongoCounter(config.General.DB, config.General.SampingTimeD())
    if err != nil {
        log.Fatalln(err)
    }

    return counter
}

func getDuration(ctx *cli.Context) time.Duration {
    if ctx.Args().Present() {
        timeSpan, err := time.ParseDuration(ctx.Args().First())
        if err != nil {
            log.Fatalln(err)
        }
        return timeSpan
    }
    return time.Hour * 1
}

func getTags(storage core.CountReader) (tagNames []string, err error) {
    tags, err := storage.Tags()
    if err != nil {
        return
    }

    tagNames = make([]string, len(tags))
    for i, tag := range tags {
        tagNames[i] = tag.Name
    }

    return
}
