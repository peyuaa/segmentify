package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/peyuaa/segmentify/data"
	"github.com/peyuaa/segmentify/handlers"

	"github.com/charmbracelet/log"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var bindAddress string = ":9090"

func main() {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "segmentify",
	})
	v := data.NewValidation()

	// create the handlers
	sh := handlers.NewSlugs(l, v)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/slugs", sh.Create)
	postR.Use(sh.MiddlewareValidateSlug)

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/slugs", sh.Get)

	// CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         bindAddress, // configure the bind address
		Handler:      ch(sm),      // set the default handler
		ErrorLog:     l.StandardLog(),
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "port", bindAddress)

		l.Fatal("Error form server", "error", s.ListenAndServe())
	}()

	// trap interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := s.Shutdown(ctx)
	if err != nil {
		l.Fatal("Error shutting down server", "error", err)
	}
}
