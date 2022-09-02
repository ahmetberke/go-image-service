package main

import (
	"github.com/ahmetberke/file-server/files"
	"github.com/ahmetberke/file-server/handlers"
	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nicholasjackson/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var address = env.String("ADDRESS", false, ":9091", "server address")
var appStatus = env.String("APP_STATUS", false, "dev", "status of the application")
var baseFilePath = env.String("FILE_PATH", false, "images", "base path of files")

func init() {
	if err := env.Parse(); err != nil {
		log.Fatalf("[error] on parsing env, %s", err.Error())
	}
}

func main() {
	localStore, err := files.NewLocalStore(*baseFilePath, files.Megabyte*5)
	if err != nil {
		log.Fatalf("[error] on creating new local store, %s", err.Error())
	}

	fileHandler := handlers.NewFilesHandler(localStore)
	gzipHandler := handlers.NewGzipHandler()

	sm := mux.NewRouter()
	corsHandler := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	postHandler := sm.Methods(http.MethodPost).Subrouter()
	postHandler.Path("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}").HandlerFunc(fileHandler.UploadREST)

	getHandler := sm.Methods(http.MethodGet).Subrouter()
	getHandler.Handle("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", http.StripPrefix("/images/", http.FileServer(http.Dir(*baseFilePath))))
	getHandler.Use(gzipHandler.GzipMiddleware)

	server := http.Server{
		Addr:         *address,
		Handler:      corsHandler(sm),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go serverRun(&server)

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Printf("Shutting down server with signal %s", sig.String())

}

func serverRun(s *http.Server) {
	log.Printf("[SERVER] Server is started at %s ", *address)
	err := s.ListenAndServe()
	if err != nil {
		log.Printf("[Error] Unable to start server: %s", err.Error())
		os.Exit(1)
	}
}
