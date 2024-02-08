package db

import (
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

type QuoteStore interface {
	ById(id uuid.UUID) (*models.Quote, error)
	Random() (*models.Quote, error)
	All() ([]*models.Quote, error)
}
