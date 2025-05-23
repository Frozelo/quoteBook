package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Frozelo/quoteBook/internal/store"
	"github.com/gorilla/mux"
)

func setupRouter() (*mux.Router, *store.Store) {
	s := store.New()
	h := &Handler{store: s}
	r := mux.NewRouter()
	r.HandleFunc("/quotes", h.GetQuotes).Methods("GET")
	r.HandleFunc("/quotes", h.PostQuote).Methods("POST")
	r.HandleFunc("/quotes/random", h.GetRandomQuote).Methods("GET")
	r.HandleFunc("/quotes/{id:[0-9]+}", h.DeleteQuote).Methods("DELETE")
	return r, s
}

func TestPostAndGetQuotes(t *testing.T) {
	r, _ := setupRouter()

	// POST /quotes
	body := []byte(`{"author":"Tester","quote":"Integration works!"}`)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.Code)
	}

	var created store.Quote
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.Author != "Tester" {
		t.Errorf("unexpected author: %s", created.Author)
	}

	// GET /quotes
	getReq := httptest.NewRequest("GET", "/quotes", nil)
	getResp := httptest.NewRecorder()
	r.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", getResp.Code)
	}
	var quotes []store.Quote
	if err := json.NewDecoder(getResp.Body).Decode(&quotes); err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Author != "Tester" {
		t.Errorf("unexpected quotes: %+v", quotes)
	}
}

func TestFilterQuotesByAuthor(t *testing.T) {
	r, s := setupRouter()
	s.Add(&store.Quote{Author: "Alice", Quote: "First"})
	s.Add(&store.Quote{Author: "Bob", Quote: "Second"})

	req := httptest.NewRequest("GET", "/quotes?author=Bob", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	var quotes []store.Quote
	if err := json.NewDecoder(resp.Body).Decode(&quotes); err != nil {
		t.Fatal(err)
	}
	if len(quotes) != 1 || quotes[0].Author != "Bob" {
		t.Errorf("expected only Bob's quotes, got: %+v", quotes)
	}
}

func TestRandomQuoteAndDelete(t *testing.T) {
	r, s := setupRouter()
	q := s.Add(&store.Quote{Author: "Test", Quote: "To be deleted"})

	req := httptest.NewRequest("GET", "/quotes/random", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.Code)
	}
	var randomQuote store.Quote
	if err := json.NewDecoder(resp.Body).Decode(&randomQuote); err != nil {
		t.Fatal(err)
	}
	if randomQuote.Id != q.Id {
		t.Errorf("random quote ID mismatch")
	}

	delReq := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(q.Id), nil)
	delResp := httptest.NewRecorder()
	r.ServeHTTP(delResp, delReq)

	if delResp.Code != http.StatusNoContent {
		t.Errorf("expected 204 No Content, got %d", delResp.Code)
	}

	delReq2 := httptest.NewRequest("DELETE", "/quotes/"+strconv.Itoa(q.Id), nil)
	delResp2 := httptest.NewRecorder()
	r.ServeHTTP(delResp2, delReq2)
	if delResp2.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d", delResp2.Code)
	}
}

func TestPostInvalidQuote(t *testing.T) {
	r, _ := setupRouter()

	body := []byte(`{}`)
	req := httptest.NewRequest("POST", "/quotes", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("expected 400 Bad Request, got %d", resp.Code)
	}
}

func TestRandomQuote_EmptyStore(t *testing.T) {
	r, _ := setupRouter()
	req := httptest.NewRequest("GET", "/quotes/random", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found for empty store, got %d", resp.Code)
	}
}
