package repository_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Linkinlog/quotes/db/dbfakes"
	"github.com/Linkinlog/quotes/models"
	"github.com/Linkinlog/quotes/repository"
	"github.com/google/uuid"
)

func TestQuote(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected *models.Quote
	}{
		"valid quote": {
			input:    "f47ac10b-58cc-0372-8567-0e02b2c3d479",
			expected: &models.Quote{Content: "This is a quote", Author: "John Doe"},
		},
		"invalid quote": {
			input:    "55555aaa-58cc-0372-8567-0e02b2c3d479",
			expected: nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeStore := &dbfakes.FakeQuoteStore{}
			fakeStore.QueryByIdReturns(tc.expected, nil)

			u, e := uuid.Parse(tc.input)
			if e != nil {
				t.Fail()
			}
			quote, err := repository.NewQuoteRepository(fakeStore).ById(u)
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(quote, tc.expected) {
				t.Errorf("expected %v; got %v", tc.expected, quote)
			}
		})
	}
}

func TestQuotes(t *testing.T) {
	tests := map[string]struct {
		expected []*models.Quote
	}{
		"some quotes": {
			expected: []*models.Quote{
				{Content: "This is a quote", Author: "John Doe"},
				{Content: "This is another quote", Author: "Jane Doe"},
			},
		},
		"no quotes": {
			expected: []*models.Quote{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeStore := &dbfakes.FakeQuoteStore{}
			fakeStore.AllReturns(tc.expected, nil)

			quotes, err := repository.NewQuoteRepository(fakeStore).All()
			if err != nil {
				t.Fail()
			}

			if !reflect.DeepEqual(quotes, tc.expected) {
				t.Errorf("expected %v; got %v", tc.expected, quotes)
			}
		})
	}
}

func TestRandomQuote(t *testing.T) {
	tests := map[string]struct {
		updateMock    func(*dbfakes.FakeQuoteStore)
		expectedQuote *models.Quote
		expectError   bool
	}{
		"valid quote": {
			expectedQuote: &models.Quote{Content: "This is a quote", Author: "John Doe"},
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.AllReturns([]*models.Quote{
					{Content: "This is a quote", Author: "John Doe"},
				}, nil)
			},
		},
		"want error": {
			expectError: true,
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.AllReturns(nil, errors.New("generic error"))
			},
		},
		"no quotes": {
			expectedQuote: &models.Quote{},
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.AllReturns([]*models.Quote{}, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeStore := &dbfakes.FakeQuoteStore{}

			tc.updateMock(fakeStore)

			quote, rErr := repository.NewQuoteRepository(fakeStore).Random()

			if rErr != nil && !tc.expectError {
				t.Fail()
			}

			if !tc.expectError && !reflect.DeepEqual(quote, tc.expectedQuote) {
				t.Errorf("expected %v; got %v", tc.expectedQuote, quote)
			} else if tc.expectError && rErr == nil {
				t.Errorf("expected error; got nil")
			}
		})
	}
}

func TestInsertQuote(t *testing.T) {
	tests := map[string]struct {
		inputQuote *models.Quote
	}{
		"valid quote": {
			inputQuote: &models.Quote{Content: "This is a quote", Author: "John Doe"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeStore := &dbfakes.FakeQuoteStore{}

			rErr := repository.NewQuoteRepository(fakeStore).Insert(tc.inputQuote.Content, tc.inputQuote.Author)

			if rErr != nil {
				t.Fail()
			}
		})
	}
}

func TestUpdateQuote(t *testing.T) {
	tests := map[string]struct {
		inputQuote *models.Quote
		wantError  bool
		updateMock func(*dbfakes.FakeQuoteStore)
	}{
		"valid quote is updated": {
			inputQuote: &models.Quote{Content: "This is an updated quote", Author: "John Doe"},
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.QueryByIdReturns(&models.Quote{Content: "This is a quote", Author: "John Doe"}, nil)
			},
		},
		"error in queryById": {
			inputQuote: &models.Quote{Content: "apple", Author: "John"},
			wantError:  true,
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.QueryByIdReturns(nil, errors.New("generic error"))
			},
		},
		"error in update": {
			inputQuote: &models.Quote{Content: "apple", Author: "John"},
			wantError:  true,
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.QueryByIdReturns(&models.Quote{Content: "This is a quote", Author: "John Doe"}, nil)
				fakeStore.UpdateReturns(errors.New("generic error"))
			},
		},
		"quote not found": {
			inputQuote: &models.Quote{Content: "apple", Author: "John"},
			updateMock: func(fakeStore *dbfakes.FakeQuoteStore) {
				fakeStore.QueryByIdReturns(nil, nil)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeStore := &dbfakes.FakeQuoteStore{}
			tc.updateMock(fakeStore)

			q, rErr := repository.NewQuoteRepository(fakeStore).Update(uuid.New(), tc.inputQuote.Content, tc.inputQuote.Author)

			if rErr != nil && !tc.wantError {
				t.Fail()
			}

			if q != nil && tc.inputQuote.Content != q.Content {
				t.Errorf("expected %v; got %v", tc.inputQuote.Content, q.Content)
			}
		})
	}
}
