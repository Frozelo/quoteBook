package store

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	appErrors "github.com/Frozelo/quoteBook/pkg/errors"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAddAndGetQuotes(t *testing.T) {
	s := New()

	s.Add(&Quote{Author: "Author1", Quote: "Quote text 1"})
	s.Add(&Quote{Author: "Author2", Quote: "Quote text 2"})

	all := s.Get()
	if len(all) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(all))
	}

	if all[0].Author != "Author1" || all[1].Author != "Author2" {
		t.Error("authors mismatch")
	}
	if all[0].Quote != "Quote text 1" || all[1].Quote != "Quote text 2" {
		t.Error("quotes mismatch")
	}

	if all[0].Id != 1 || all[1].Id != 2 {
		t.Error("ID mismatch")
	}

	byAuthor := s.GetByAuthor("Author1")
	if len(byAuthor) != 1 || byAuthor[0].Author != "Author1" {
		t.Error("author filter broken")
	}

	empty := s.GetByAuthor("Unknown")
	if len(empty) != 0 {
		t.Error("expected empty slice for unknown author")
	}
}

func TestGetRandom(t *testing.T) {
	s := New()
	_, err := s.GetRandom()
	if !errors.Is(err, appErrors.ErrNoQuotes) {
		t.Error("expected ErrNoQuotes on empty store")
	}

	s.Add(&Quote{Author: "A", Quote: "Q1"})
	q, err := s.GetRandom()
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if q.Author != "A" || q.Quote != "Q1" {
		t.Error("random quote mismatch")
	}
}

func TestDelete(t *testing.T) {
	s := New()
	q := s.Add(&Quote{Author: "A", Quote: "Q1"})
	err := s.Delete(q.Id)
	if err != nil {
		t.Error("expected delete to succeed, got:", err)
	}

	err = s.Delete(q.Id)
	if !errors.Is(err, appErrors.ErrQuoteNotFound) {
		t.Error("expected ErrQuoteNotFound, got:", err)
	}
	if len(s.Get()) != 0 {
		t.Error("expected store to be empty after delete")
	}
}

func TestGetRandom_EmptyStore(t *testing.T) {
	s := New()
	_, err := s.GetRandom()
	if !errors.Is(err, appErrors.ErrNoQuotes) {
		t.Errorf("expected ErrNoQuotes, got %v", err)
	}
}

func TestDelete_NotFound(t *testing.T) {
	s := New()
	err := s.Delete(999)
	if !errors.Is(err, appErrors.ErrQuoteNotFound) {
		t.Errorf("expected ErrQuoteNotFound, got %v", err)
	}
}

func TestGetByAuthor_UnknownAuthor(t *testing.T) {
	s := New()
	s.Add(&Quote{Author: "A", Quote: "Q1"})
	quotes := s.GetByAuthor("NoSuchAuthor")
	if len(quotes) != 0 {
		t.Errorf("expected 0 quotes, got %d", len(quotes))
	}
}

func TestDelete_EmptyStore(t *testing.T) {
	s := New()
	err := s.Delete(1)
	if !errors.Is(err, appErrors.ErrQuoteNotFound) {
		t.Errorf("expected ErrQuoteNotFound when deleting from empty store, got %v", err)
	}
}
