package store

import (
	"math/rand"
	"slices"
	"sync"

	appErrors "github.com/Frozelo/quoteBook/pkg/errors"
)

type Quote struct {
	Id     int    `json:"id"`
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type Store struct {
	mu     sync.Mutex
	quotes []Quote
	nextId int
}

func New() *Store {
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

func (qs *Store) Get() []Quote {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	return qs.quotes
}

func (qs *Store) GetByAuthor(author string) []Quote {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	filtered := make([]Quote, 0)

	for _, q := range qs.quotes {
		if q.Author == author {
			filtered = append(filtered, q)
		}
	}

	if len(filtered) == 0 {
		return filtered
	}

	return filtered
}

func (qs *Store) GetRandom() (Quote, error) {
	qs.mu.Lock()
	defer qs.mu.Unlock()

	if len(qs.quotes) == 0 {
		return Quote{}, appErrors.ErrNoQuotes
	}

	idx := rand.Intn(len(qs.quotes))

	return qs.quotes[idx], nil
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
	return appErrors.ErrQuoteNotFound
}
