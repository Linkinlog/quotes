package repository

import (
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

func (r *QuoteRepository) Quote(id uuid.UUID) (*models.Quote, error) {
	return r.store.ById(id)
}

func (r *QuoteRepository) Quotes() ([]*models.Quote, error) {
	return r.store.All()
}

func (r *QuoteRepository) RandomQuote() (*models.Quote, error) {
	return r.store.Random()
}

func MakeQuotes() []*models.Quote {
	// TODO - move this to a file
	return []*models.Quote{
		models.NewQuote("Hello, World.", "Log"),
		models.NewQuote("Heaven or hell, love or hate, No matter where I turn I meet myself. Holding life precious is Just living with all intensity Holding life precious.", "Opening the Hand of Thought p.81- Kōshō Uchiyama"),
		models.NewQuote("Regardless of which component you begin at, it is impossible to follow the dependency relationships and wind up back at that component. This structure has no cycles. It is a directed acyclic graph(DAG)", "Clean Architecture Ch 14 - Robert C. Martin"),
		models.NewQuote("The key point here is our programmers are Googlers, they’re not researchers. They’re typically, fairly young, fresh out of school, probably learned Java, maybe learned C or C++, probably learned Python. They’re not capable of understanding a brilliant language but we want to use them to build good software. So, the language that we give them has to be easy for them to understand and easy to adopt.", "Rob Pike"),
	}
}
