package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"

	"deck-service/handlers"
	"deck-service/storage"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func main() {
	var address string
	flag.StringVar(&address, "address", ":8089", "http address")
	flag.Parse()

	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync() // flushes buffer, if any

	// init storage
	storage := storage.NewConcurrentMemory()

	r := mux.NewRouter()
	handlers := handlers.New(storage, logger)

	r.HandleFunc("/deck", handlers.CreateDeckHandler).Methods(http.MethodPost)
	r.HandleFunc("/decks/{id}", handlers.OpenDeckHandler).Methods(http.MethodGet)
	r.HandleFunc("/decks/{id}/draw", handlers.DrawCard).Methods(http.MethodPost)
	go logger.Fatal("failed to serve", zap.Error(http.ListenAndServe(address, r)))
	<-shutdown()
}

// shutdown listens SIGINT signal and shutdown application.
func shutdown() <-chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	done := make(chan struct{})
	go func() {
		for range c {
			close(done)
		}
	}()
	return done
}
