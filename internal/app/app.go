package app

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Frozelo/quoteBook/internal/handlers"
	"github.com/Frozelo/quoteBook/internal/middelware"
	"github.com/Frozelo/quoteBook/internal/server"
	"github.com/Frozelo/quoteBook/internal/store"
	"github.com/gorilla/mux"
)

func Run() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))

	router := mux.NewRouter()
	quoteStore := store.New()
	quoteHandler := handlers.New(quoteStore)

	router.HandleFunc("/quotes", quoteHandler.GetQuotes).Methods("GET")
	router.HandleFunc("/quotes", quoteHandler.PostQuote).Methods("POST")
	router.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.DeleteQuote).Methods("DELETE")

	router.Use(middelware.LoggingMiddleware(logger))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Info("starting new servet at port", "port", "8080")
	httpServer := server.New(router)

	select {
	case s := <-interrupt:
		log.Printf("app - Run - signal %s", s.String())

	case err := <-httpServer.Notify():
		log.Printf("app - Run - httpServer.Notify: %v", err)

		err = httpServer.Shutdown()
		if err != nil {
			log.Printf("app - Run - httpServer.Shutdown: %v", err)
		}
	}

}
