package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"slices"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

type AddQuoteReq struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type Quote struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

var ErrQuoteNotFound = errors.New("quote not found")

type Store struct {
	mu     sync.Mutex
	quotes []Quote
	nextId int
}

func NewStore() *Store {
	return &Store{
		quotes: make([]Quote, 0),
		nextId: 1,
	}
}

func (qs *Store) Add(quote *Quote) Quote {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	newQuote := Quote{
		Id:     qs.nextId,
		Author: quote.Author,
		Quote:  quote.Quote,
	}

	qs.quotes = append(
		qs.quotes,
		newQuote,
	)

	qs.nextId++
	return newQuote
}

func (qs *Store) Get() ([]Quote, error) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	return qs.quotes, nil
}

func (qs *Store) Delete(id int) error {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	for i, q := range qs.quotes {
		if q.Id == id {
			qs.quotes = slices.Delete(qs.quotes, i, i+1)
			return nil
		}
	}
	return ErrQuoteNotFound
}

func main() {
	router := mux.NewRouter()

	quoteStore := NewStore()

	router.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		var req AddQuoteReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newQuote := quoteStore.Add(&Quote{
			Author: req.Author,
			Quote:  req.Quote,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newQuote)
	}).Methods("POST")

	router.HandleFunc("/quotes", func(w http.ResponseWriter, r *http.Request) {
		quotes, err := quoteStore.Get()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(quotes)
	}).Methods("GET")

	router.HandleFunc("/quotes/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := quoteStore.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}).Methods("DELETE")

	router.HandleFunc("/heartbeat", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	log.Println("Server started on port 8080")
	http.ListenAndServe(":8080", router)
}
