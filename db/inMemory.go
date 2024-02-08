package db

import (
	"math/rand"

	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

type inMemoryStore struct {
	quotes []*models.Quote
}

func NewInMemoryStore(quotes []*models.Quote) QuoteStore {
	return &inMemoryStore{quotes: quotes}
}

func (s *inMemoryStore) AddQuote(q *models.Quote) error {
	s.quotes = append(s.quotes, q)
	return nil
}

func (s *inMemoryStore) All() ([]*models.Quote, error) {
	return s.quotes, nil
}

func (s *inMemoryStore) ById(id uuid.UUID) (*models.Quote, error) {
	for _, q := range s.quotes {
		if q.Id() == id {
			return q, nil
		}
	}
	return &models.Quote{}, nil
}

func (s *inMemoryStore) Random() (*models.Quote, error) {
	randIndex := rand.Intn(len(s.quotes))
	return s.quotes[randIndex], nil
}
