package main

import (
	"github.com/hashtock/tracker/listener"
	"log"
	"net/http"

	authClient "github.com/hashtock/auth/client"

	"github.com/hashtock/tracker/conf"
	"github.com/hashtock/tracker/core"
	"github.com/hashtock/tracker/storage"
	"github.com/hashtock/tracker/webapi"
)

func main() {
	cfg := conf.GetConfig()
	cfg.General.Timeout = 0 // No timeout

	counter, err := storage.NewMongoCounter(cfg.General.DB, cfg.General.SampingTime)
	if err != nil {
		log.Fatalln(err)
	}

	tagNames, err := core.TagNames(counter)
	if err != nil {
		log.Fatalln(err)
	}

	listenerOptions := listener.Options{
		TagUpdateTime: cfg.General.TagUpdateTime,
		SampingTime:   cfg.General.SampingTime,
		Verbose:       true,
	}
	twitterListener := listener.NewTwitterListener(tagNames, cfg.General.Timeout, cfg.General.UpdateTime, cfg.Auth)
	go listener.Listen(twitterListener, counter, listenerOptions)

	handlerOptions := webapi.Options{
		Serializer: new(webapi.WebAPISerializer),
		Counter:    counter,
		WhoClient:  authClient.NewClient("http://localhost:4000/auth/who/"),
	}
	handler := webapi.Handlers(handlerOptions)
	err = http.ListenAndServe(cfg.General.ServeAddress, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
