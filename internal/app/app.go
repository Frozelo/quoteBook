package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	quoteHandler := handlers.New(logger, quoteStore)

	router.HandleFunc("/quotes", quoteHandler.GetQuotes).Methods("GET")
	router.HandleFunc("/quotes", quoteHandler.PostQuote).Methods("POST")
	router.HandleFunc("/quotes/random", quoteHandler.GetRandomQuote).Methods("GET")
	router.HandleFunc("/quotes/{id:[0-9]+}", quoteHandler.DeleteQuote).Methods("DELETE")

	router.Use(middelware.LoggingMiddleware(logger))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	logger.Info("starting new servet at port", "port", "8080")
	httpServer := server.New(router)

	if os.Getenv("SELF_CHECK") == "1" {
		go func() {
			time.Sleep(1 * time.Second)
			fmt.Println("Запуск self-check через /quotes")
			err := runSelfCheck("http://localhost:8080")
			if err != nil {
				fmt.Println("SELF_CHECK failed:", err)
				os.Exit(2)
			} else {
				fmt.Println("SELF_CHECK successful!")
				os.Exit(0)
			}
		}()
	}

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

func runSelfCheck(baseURL string) error {
	fmt.Println("== Self-check ==")
	client := &http.Client{}

	fmt.Println("\nДобавление новой цитаты")
	postBody := map[string]string{"author": "Confucius", "quote": "Life is simple, but we insist on making it complicated."}
	bodyBytes, _ := json.Marshal(postBody)
	resp, err := client.Post(baseURL+"/quotes", "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	printResp(resp)

	fmt.Println("\nПолучение всех цитат")
	resp, err = client.Get(baseURL + "/quotes")
	if err != nil {
		return err
	}
	printResp(resp)

	fmt.Println("\nПолучение случайной цитаты")
	resp, err = client.Get(baseURL + "/quotes/random")
	if err != nil {
		return err
	}
	printResp(resp)

	fmt.Println("\nФильтрация по автору")
	resp, err = client.Get(baseURL + "/quotes?author=Confucius")
	if err != nil {
		return err
	}
	printResp(resp)

	fmt.Println("\nУдаление цитаты по ID")
	req, _ := http.NewRequest("DELETE", baseURL+"/quotes/1", nil)
	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	printResp(resp)

	return nil
}

func printResp(resp *http.Response) {
	defer resp.Body.Close()
	fmt.Printf("Status: %s\n", resp.Status)
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Body:", string(body))
}
