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
	quotes, err := r.store.All()
	if err != nil {
		return nil, err
	}
	if len(quotes) == 0 {
		return &models.Quote{}, nil
	}
	randIndex := rand.Intn(len(quotes))
	return quotes[randIndex], nil
}

func (r *QuoteRepository) Insert(content, author string) error {
	quote := models.NewQuote(content, author)
	return r.store.Insert(quote)
}

func (r *QuoteRepository) Update(id uuid.UUID, content, author string) (*models.Quote, error) {
	quote, err := r.store.QueryById(id)
	if err != nil {
		return nil, err
	}
	if quote == nil {
		return nil, nil
	}
	quote.Content = content
	quote.Author = author
	err = r.store.Update(quote)
	if err != nil {
		return nil, err
	}
	return quote, nil
}
