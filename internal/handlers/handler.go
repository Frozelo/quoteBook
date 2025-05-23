package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Frozelo/quoteBook/internal/store"
	"github.com/gorilla/mux"
)

type Handler struct {
	logger *slog.Logger
	store  *store.Store
}

func New(logger *slog.Logger, store *store.Store) *Handler {
	return &Handler{logger: logger, store: store}
}

// POST /quotes
func (h *Handler) PostQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Author string `json:"author"`
		Quote  string `json:"quote"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Author == "" || req.Quote == "" {
		h.logger.Warn("missing author or quote in request", "request", req)
		http.Error(w, "author and quote required", http.StatusBadRequest)
		return
	}

	newQuote := &store.Quote{
		Author: req.Author,
		Quote:  req.Quote,
	}

	q := h.store.Add(newQuote)
	h.logger.Info("quote added", "author", q.Author, "id", q.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(q)
}

// GET /quotes
func (h *Handler) GetQuotes(w http.ResponseWriter, r *http.Request) {
	author := r.URL.Query().Get("author")
	w.Header().Set("Content-Type", "application/json")
	if author == "" {
		h.logger.Info("fetch all quotes")
		json.NewEncoder(w).Encode(h.store.Get())
		return
	}
	h.logger.Info("fetch quotes by author", "author", author)
	quotes := h.store.GetByAuthor(author)
	json.NewEncoder(w).Encode(quotes)
}

// GET /quotes/random
func (h *Handler) GetRandomQuote(w http.ResponseWriter, r *http.Request) {
	quote, err := h.store.GetRandom()
	if err != nil {
		h.logger.Warn("no quotes for random selection")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	h.logger.Info("random quote fetched", "id", quote.Id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quote)
}

// DELETE /quotes/{id}
func (h *Handler) DeleteQuote(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("invalid id for delete", "id", idStr)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = h.store.Delete(id)
	if err != nil {
		h.logger.Warn("delete failed, quote not found", "id", id)
		http.Error(w, "quote not found", http.StatusNotFound)
		return
	}
	h.logger.Info("quote deleted", "id", id)
	w.WriteHeader(http.StatusNoContent)
}
