package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/peyuaa/segmentify/data"
	"github.com/peyuaa/segmentify/db"
	"github.com/peyuaa/segmentify/handlers"

	"github.com/charmbracelet/log"
	"github.com/go-openapi/runtime/middleware"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	// DBConnectionString is a name of the environment variable
	//that contains the connection string to the database
	DBConnectionString = "DB_CONNECTION_STRING"
)

var bindAddress = ":9090"

func main() {
	l := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Prefix:          "segmentify",
	})
	v := data.NewValidation()

	// get the environment variables
	dbConnectionString := os.Getenv(DBConnectionString)
	if dbConnectionString == "" {
		l.Fatal("DB_CONNECTION_STRING isn't set")
	}

	l.Info("Connecting to postgresql database")

	l.Info("Waiting for postgresql database to start")
	time.Sleep(10 * time.Second)

	// set up the database connection
	dbConn, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		l.Fatal("Unable to connect to database", "error", err)
	}
	defer func() {
		err := dbConn.Close()
		if err != nil {
			l.Error("Unable to close database connection", "error", err)
		}
	}()

	// establish connection to the database
	err = dbConn.Ping()
	if err != nil {
		l.Fatal("Unable to ping database", "error", err)
	}

	l.Info("Connected to postgresql database")

	// create postgresql wrapper
	dbWrap := db.New(l, dbConn)

	// create new database struct
	segmentifyDB := data.New(l, dbWrap)

	// create the handlers
	sh := handlers.NewSegments(l, v, segmentifyDB)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	postR := sm.Methods(http.MethodPost).Subrouter()
	segR := postR.Path("/segments").Subrouter()
	segR.HandleFunc("", sh.CreateSegment)
	segR.Use(sh.MiddlewareValidateSegment)

	userR := postR.Path("/segments/users").Subrouter()
	userR.HandleFunc("", sh.ChangeUsersSegments)
	userR.Use(sh.MiddlewareValidateUser)

	getR := sm.Methods(http.MethodGet).Subrouter()
	// serve directory with user history files
	getR.PathPrefix("/history/").Handler(http.StripPrefix("/history/", http.FileServer(http.Dir("history"))))

	getR.HandleFunc("/segments", sh.GetSegments)
	getR.HandleFunc("/segments/{slug:[a-zA-Z_0-9]+}", sh.GetBySlug)
	getR.HandleFunc("/segments/users/{id:[0-9]+}", sh.GetActiveSegments)
	getR.HandleFunc("/segments/users/{id:[0-9]+}/history", sh.UserHistory)

	// handlers for documentation
	opts := middleware.RedocOpts{
		SpecURL: "/swagger.yaml",
	}
	dh := middleware.Redoc(opts, nil)
	getR.Handle("/docs", dh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/segments/{slug:[a-zA-Z_0-9]+}", sh.Delete)

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

	err = s.Shutdown(ctx)
	if err != nil {
		l.Fatal("Error shutting down server", "error", err)
	}
}
