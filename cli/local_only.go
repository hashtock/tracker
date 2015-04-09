package cli

import (
	"fmt"
	"log"
	"reflect"
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

	exitSync := make(chan struct{})

	tagNames, err := getTags(counter)
	if err != nil {
		log.Fatalln(err)
	}

	twitterListener := listener.NewTwitterListener(tagNames, cfg.General.Timeout, cfg.General.UpdateTime, cfg.Auth)
	countCh := twitterListener.Listen()

	go func() {
		watcher := time.NewTicker(cfg.General.TagUpdateTime)
		defer watcher.Stop()

		for {
			select {
			case <-watcher.C:
				newTags, err := getTags(counter)
				if err != nil {
					log.Println("Could not get new tags:", err)
				} else if !reflect.DeepEqual(newTags, twitterListener.Tags()) {
					log.Println("Setting new list of tags to track")
					twitterListener.SetTags(newTags)
				}
			case <-exitSync:
				twitterListener.Stop()
				return
			}
		}
	}()

	for countMap := range countCh {
		now := time.Now().Truncate(cfg.General.SampingTime)
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

	close(exitSync)
}

func cmdWebAPI(ctx *cli.Context) {
	counter := getCounterRW(ctx)
	cfg := conf.GetConfig()
	cfg.General.Timeout = 0 // No timeout
	go cmdListen(ctx)

	webapi.RunWebAPI(counter)
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

func cmdPrintConfHelp(ctx *cli.Context) {
	conf.PrintConfHelp()
}
