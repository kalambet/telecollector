package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/kalambet/telecollector/store"

	"github.com/kalambet/telecollector/telecollector"
)

var (
	ErrTGTokenEmpty = errors.New("server: telegram token is empty")
)

type server struct {
	port        int
	router      *http.ServeMux
	msgService  telecollector.MessageService
	credService telecollector.CredentialService
}

type response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewServer(ms telecollector.MessageService, cred telecollector.CredentialService) (*server, error) {
	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, err
	}

	res := &server{
		port:        port,
		msgService:  ms,
		credService: cred,
		router:      http.NewServeMux(),
	}

	token := os.Getenv("TG_TOKEN")
	if len(token) == 0 {
		return nil, ErrTGTokenEmpty
	}
	res.routes(token)

	return res, nil
}

func (s *server) StartServer() {
	srv := http.Server{
		Handler: s.router,
		Addr:    fmt.Sprintf(":%d", s.port),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("server: execution was interrupted: %s\n", err.Error())
		}
	}()

	// Setting up signal capturing
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	// Waiting for SIGINT (pkill -2)
	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer s.stopServer()
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("server: error while shutdown: %s\n", err.Error())
	}
}

func (s *server) stopServer() {
	err := store.Shutdown()
	if err != nil {
		log.Fatalf("server: store shutdown error: %s", err.Error())
	}
}

func (s *server) respond(w http.ResponseWriter, st int, m string) {

	w.Header().Add("content-type", "application/json")
	w.Header().Add("accept", "application/json")

	w.WriteHeader(http.StatusOK)

	data := response{
		Status:  st,
		Message: m,
	}

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		log.Printf("server: error creating response: %s", err.Error())
	}
}

func (s *server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, http.StatusOK, "All your base are belong to us!")
	}
}
