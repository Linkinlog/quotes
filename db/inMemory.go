package db

import (
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

type inMemoryStore struct {
	quotes []*models.Quote
}

func NewInMemoryStore(quotes []*models.Quote) QuoteStore {
	return &inMemoryStore{quotes: quotes}
}

func (s *inMemoryStore) Insert(q *models.Quote) error {
	s.quotes = append(s.quotes, q)
	return nil
}

func (s *inMemoryStore) All() ([]*models.Quote, error) {
	return s.quotes, nil
}

func (s *inMemoryStore) QueryById(id uuid.UUID) (*models.Quote, error) {
	for _, q := range s.quotes {
		if q.Id == id {
			return q, nil
		}
	}
	return &models.Quote{}, nil
}

func (s *inMemoryStore) Update(q *models.Quote) error {
	for i, quote := range s.quotes {
		if quote.Id == q.Id {
			s.quotes[i] = q
			return nil
		}
	}
	return nil
}

func (s *inMemoryStore) Delete(id uuid.UUID) error {
	for i, q := range s.quotes {
		if q.Id == id {
			s.quotes[i].Disapprove()
			return nil
		}
	}
	return nil
}
