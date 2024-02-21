package db_test

import (
	"testing"

	"github.com/Linkinlog/quotes/db"
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
)

func exampleQuotes() []*models.Quote {
	u1, _ := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d479")
	u2, _ := uuid.Parse("f47ac10b-58cc-0372-8567-0e02b2c3d480")
	return []*models.Quote{
		{Id: u1, Content: "This is a quote", Author: "John Doe"},
		{Id: u2, Content: "This is another quote", Author: "Jane Doe"},
	}
}

func TestInMemoryStore_All(t *testing.T) {
	tests := map[string]struct {
		quotes   []*models.Quote
		expected []*models.Quote
	}{
		"some quotes": {
			quotes:   exampleQuotes(),
			expected: exampleQuotes(),
		},
		"no quotes": {
			quotes:   []*models.Quote{},
			expected: []*models.Quote{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			store := db.NewInMemoryStore(tc.quotes)
			quotes, err := store.All()
			if err != nil {
				t.Fail()
			}

			if len(quotes) != len(tc.expected) {
				t.Errorf("expected %v; got %v", tc.expected, quotes)
			}
		})
	}
}

func TestInMemoryStore_QueryById(t *testing.T) {
	tests := map[string]struct {
		quotes   []*models.Quote
		id       string
		expected *models.Quote
	}{
		"valid quote": {
			quotes:   exampleQuotes(),
			id:       "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: exampleQuotes()[0],
		},
		"invalid quote": {
			quotes:   exampleQuotes(),
			id:       "55555555-58cc-0372-8567-0e02b2c3d480",
			expected: &models.Quote{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			store := db.NewInMemoryStore(tc.quotes)
			id, _ := uuid.Parse(tc.id)
			quote, err := store.QueryById(id)
			if err != nil {
				t.Fail()
			}

			if quote.Id != tc.expected.Id {
				t.Errorf("expected %v; got %v", tc.expected, quote)
			}
		})
	}
}

func TestInMemoryStore_Insert(t *testing.T) {
	tests := map[string]struct {
		quotes   []*models.Quote
		input    *models.Quote
		expected []*models.Quote
	}{
		"valid quote": {
			quotes:   exampleQuotes(),
			input:    &models.Quote{Content: "This is a new quote", Author: "John Doe"},
			expected: append(exampleQuotes(), &models.Quote{Content: "This is a new quote", Author: "John Doe"}),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			store := db.NewInMemoryStore(tc.quotes)
			err := store.Insert(tc.input)
			if err != nil {
				t.Fail()
			}

			quotes, _ := store.All()
			if len(quotes) != len(tc.expected) {
				t.Errorf("expected %v; got %v", tc.expected, quotes)
			}
		})
	}
}
