package main

import (
	"log"

	"github.com/kalambet/telecollector/telecollector"

	"github.com/kalambet/telecollector/http"

	"github.com/kalambet/telecollector/store"
)

var msg telecollector.MessageService

func init() {
	var err error
	msg, err = store.NewMessagesService()
	if err != nil {
		log.Fatalf("stratup: error creating messaging service: %s", err.Error())
	}
}

func main() {
	srv, err := http.NewServer(msg)
	if err != nil {
		log.Fatalf("startup: error initializing server: %s", err.Error())
	}

	srv.StartServer()
}
