package handlers_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Linkinlog/quotes/db/dbfakes"
	"github.com/Linkinlog/quotes/handlers"
	"github.com/Linkinlog/quotes/models"
	"github.com/Linkinlog/quotes/repository"
	"github.com/google/uuid"
)

func TestHandler_HandleRoutes(t *testing.T) {
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
	}{
		"valid repository": {
			quoteRepository: &repository.QuoteRepository{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fail()
			}

			_ = handler.HandleRoutes("/")
		})
	}
}

func TestHandler_HandleLanding(t *testing.T) {
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
		forceError      bool
	}{
		"valid repository": {
			quoteRepository: &repository.QuoteRepository{},
		},
		"invalid repository": {
			quoteRepository: &repository.QuoteRepository{},
			forceError:      true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fail()
			}

			req := httptest.NewRequest("GET", "/", nil)
			rr := httptest.NewRecorder()

			handler.HandleLanding(rr, req)

			if rr.Code != 200 {
				t.Fail()
			}
		})
	}
}

func TestHandler_HandleRandomQuote(t *testing.T) {
	fakeDb := &dbfakes.FakeQuoteStore{}
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
		wantErr         bool
	}{
		"valid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
		},
		"invalid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.AllReturns(nil, errors.New("error"))
			} else {
				fakeDb.AllReturns([]*models.Quote{models.NewQuote("foo", "bar")}, nil)
			}
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fatal("handler is nil")
			}

			req := httptest.NewRequest("GET", "/random", nil)
			rr := httptest.NewRecorder()

			handler.HandleRandomQuote(rr, req)

			if rr.Code != 200 && !tc.wantErr {
				t.Fail()
			} else if rr.Code != 500 && tc.wantErr {
				t.Fail()
			}
		})
	}
}

func TestHandler_HandleQuoteById(t *testing.T) {
	fakeDb := &dbfakes.FakeQuoteStore{}
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
		wantErr         bool
		wantBadRequest  bool
		id              uuid.UUID
	}{
		"valid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
			id:              uuid.New(),
		},
		"invalid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
			wantErr:         true,
			id:              uuid.New(),
		},
		"invalid id": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
			wantBadRequest:  true,
			id:              uuid.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.QueryByIdReturns(nil, errors.New("error"))
			} else {
				fakeDb.QueryByIdReturns(&models.Quote{Id: tc.id, Content: "hi", Author: "me"}, nil)
			}
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fatal("handler is nil")
			}

			path := "/quotes/" + tc.id.String()
			req := httptest.NewRequest("GET", path, nil)
			rr := httptest.NewRecorder()
			req.SetPathValue("id", tc.id.String())

			handler.HandleQuoteById(rr, req)

			if rr.Code != 200 && !(tc.wantErr || tc.wantBadRequest) {
				t.Fatal("expected 200, got", rr.Code)
			} else if rr.Code != 500 && tc.wantErr {
				t.Fatal("expected 500, got", rr.Code)
			} else if rr.Code != 400 && tc.wantBadRequest {
				t.Fatal("expected 400, got", rr.Code)
			}
		})
	}
}

func TestHandler_HandleAbout(t *testing.T) {
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
	}{
		"valid repository": {
			quoteRepository: &repository.QuoteRepository{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fail()
			}

			req := httptest.NewRequest("GET", "/about", nil)
			rr := httptest.NewRecorder()

			handler.HandleAbout(rr, req)

			if rr.Code != 200 {
				t.Fail()
			}
		})
	}
}

func TestHandler_HandleQuotes(t *testing.T) {
	fakeDb := &dbfakes.FakeQuoteStore{}
	tests := map[string]struct {
		quoteRepository *repository.QuoteRepository
		wantErr         bool
	}{
		"valid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
		},
		"invalid repository": {
			quoteRepository: repository.NewQuoteRepository(fakeDb),
			wantErr:         true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.AllReturns(nil, errors.New("error"))
			} else {
				fakeDb.AllReturns([]*models.Quote{models.NewQuote("foo", "bar")}, nil)
			}
			handler := handlers.NewHandler(tc.quoteRepository)
			if handler == nil {
				t.Fatal("handler is nil")
			}

			req := httptest.NewRequest("GET", "/quotes", nil)
			rr := httptest.NewRecorder()

			handler.HandleQuotes(rr, req)

			if rr.Code != 200 && !tc.wantErr {
				t.Fail()
			} else if rr.Code != 500 && tc.wantErr {
				t.Fail()
			}
		})
	}
}
