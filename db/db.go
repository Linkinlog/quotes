package db

import (
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . QuoteStore

type QuoteStore interface {
	Insert(*models.Quote) error
	QueryById(uuid.UUID) (*models.Quote, error)
	All() ([]*models.Quote, error)
	Update(*models.Quote) error
	Delete(uuid.UUID) error
}
