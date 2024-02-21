package repository

import (
	"math/rand"

	"github.com/Linkinlog/quotes/db"
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

type QuoteRepository struct {
	store db.QuoteStore
}

func NewQuoteRepository(store db.QuoteStore) *QuoteRepository {
	return &QuoteRepository{store: store}
}

func (r *QuoteRepository) ById(id uuid.UUID) (*models.Quote, error) {
	return r.store.QueryById(id)
}

func (r *QuoteRepository) All() ([]*models.Quote, error) {
	return r.store.All()
}

func (r *QuoteRepository) Random() (*models.Quote, error) {
	quotes, err := r.All()
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return &models.Quote{}, nil
	}
	randIndex := rand.Intn(len(quotes))
	return quotes[randIndex], nil
}
