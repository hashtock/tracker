package main

import (
	"log"
	"net/http"

	authClient "github.com/hashtock/auth/client"
	"github.com/nats-io/nats"

	"github.com/hashtock/tracker/conf"
	"github.com/hashtock/tracker/core"
	"github.com/hashtock/tracker/listener"
	"github.com/hashtock/tracker/storage"
	"github.com/hashtock/tracker/webapi"
)

func main() {
	cfg := conf.GetConfig()
	cfg.General.Timeout = 0 // No timeout

	nc, err := nats.Connect(cfg.General.NATS)
	if err != nil {
		log.Fatalln(err)
	}
	msgConnection, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		log.Fatalln(err)
	}
	defer msgConnection.Close()

	counter, err := storage.NewMongoCounter(cfg.General.DB, cfg.General.SampingTime)
	if err != nil {
		log.Fatalln(err)
	}

	tagNames, err := core.TagNames(counter)
	if err != nil {
		log.Fatalln(err)
	}

	noticiator := NewMsgNoticiator(counter, msgConnection)
	listenerOptions := listener.Options{
		TagUpdateTime: cfg.General.TagUpdateTime,
		SampingTime:   cfg.General.SampingTime,
		Verbose:       true,
		Notificator:   noticiator,
	}
	twitterListener := listener.NewTwitterListener(tagNames, cfg.General.Timeout, cfg.General.UpdateTime, cfg.Auth)
	go listener.Listen(twitterListener, counter, listenerOptions)

	whoClient, whoErr := authClient.NewClient(cfg.General.AuthAddress)
	if whoErr != nil {
		log.Fatalln(err)
	}
	handlerOptions := webapi.Options{
		Serializer: new(webapi.WebAPISerializer),
		Counter:    counter,
		WhoClient:  whoClient,
	}
	handler := webapi.Handlers(handlerOptions)
	err = http.ListenAndServe(cfg.General.ServeAddress, handler)
	if err != nil {
		log.Fatalln(err)
	}
}
