package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Frozelo/quoteBook/internal/store"
	"github.com/gorilla/mux"
)

type Handler struct {
	store *store.Store
}

func New(store *store.Store) *Handler {
	return &Handler{store: store}
}

// POST /quotes
func (h *Handler) PostQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Author string `json:"author"`
		Quote  string `json:"quote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Author == "" || req.Quote == "" {
		http.Error(w, "author and quote required", http.StatusBadRequest)
		return
	}

	newQuote := &store.Quote{
		Author: req.Author,
		Quote:  req.Quote,
	}

	q := h.store.Add(newQuote)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// GET /quotes
func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	w.Header().Set("Content-Type", "application/json")
	if author == "" {
		json.NewEncoder(w).Encode(h.store.Get())
		return
	}
	quotes := h.store.GetByAuthor(author)
	json.NewEncoder(w).Encode(quotes)
}

// GET /quotes/random
func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.store.GetRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// DELETE /quotes/{id}
func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = h.store.Delete(id)
	if err != nil {
		http.Error(w, "quote not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
