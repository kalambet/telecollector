package main

import (
	"log"

	"github.com/kalambet/telecollector/telecollector"

	"github.com/kalambet/telecollector/http"

	"github.com/kalambet/telecollector/store"
)

var msg telecollector.MessageService
var cred telecollector.CredentialService

func init() {
	var err error
	msg, err = store.NewMessagesService()
	if err != nil {
		log.Fatalf("stratup: error initializing messaging service: %s", err.Error())
	}

	cred, err = store.NewCrenetialService()
	if err != nil {
		log.Fatalf("stratup: error initializing credential service: %s", err.Error())
	}
}

func main() {
	srv, err := http.NewServer(msg, cred)
	if err != nil {
		log.Fatalf("startup: error initializing server: %s", err.Error())
	}

	srv.StartServer()
}
