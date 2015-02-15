package cli

import (
    "log"
    "time"

    "github.com/codegangsta/cli"

    "github.com/hashtock/tracker/client"
    "github.com/hashtock/tracker/conf"
)

func getRemoteClient(ctx *cli.Context) *client.Tracker {
    if remote := ctx.GlobalString("remote"); remote != "" {
        remoteConfig := conf.GetRemoteConfig(remote)
        return client.NewTracker(remoteConfig)
    }

    return nil
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
