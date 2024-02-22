package handlers_test

import (
	"errors"
	"net/http"
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
			handler := handlers.NewHandler(tc.quoteRepository, "secret", true)
			if handler == nil {
				t.Fail()
			}

			_ = handler.HandleRoutes("/")
		})
	}
}

func TestHandler_HandleLanding(t *testing.T) {
	tests := map[string]struct{}{
		"valid repository": {},
	}

	for name := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(&repository.QuoteRepository{}, "secret", true)
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
		wantErr bool
	}{
		"valid repository": {},
		"invalid repository": {
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.AllReturns(nil, errors.New("error"))
			} else {
				fakeDb.AllReturns([]*models.Quote{models.NewQuote("foo", "bar")}, nil)
			}
			handler := handlers.NewHandler(repository.NewQuoteRepository(fakeDb), "secret", true)
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
		wantErr        bool
		wantBadRequest bool
		id             uuid.UUID
	}{
		"valid repository": {
			id: uuid.New(),
		},
		"invalid repository": {
			wantErr: true,
			id:      uuid.New(),
		},
		"invalid id": {
			wantBadRequest: true,
			id:             uuid.Nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.QueryByIdReturns(nil, errors.New("error"))
			} else {
				fakeDb.QueryByIdReturns(&models.Quote{Id: tc.id, Content: "hi", Author: "me"}, nil)
			}
			handler := handlers.NewHandler(repository.NewQuoteRepository(fakeDb), "secret", true)
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
	tests := map[string]struct{}{
		"valid": {},
	}

	for name := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(&repository.QuoteRepository{}, "secret", true)
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
		wantErr bool
	}{
		"valid repository": {},
		"invalid repository": {
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr {
				fakeDb.AllReturns(nil, errors.New("error"))
			} else {
				fakeDb.AllReturns([]*models.Quote{models.NewQuote("foo", "bar")}, nil)
			}
			handler := handlers.NewHandler(repository.NewQuoteRepository(fakeDb), "secret", true)
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

func TestHandler_HandleAuthMiddleware(t *testing.T) {
	tests := map[string]struct {
		shouldAuth bool
	}{
		"authenticated": {
			shouldAuth: true,
		},
		"unauthenticated": {
			shouldAuth: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(&repository.QuoteRepository{}, "secret", true)
			if handler == nil {
				t.Fail()
			}

			rr := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)

			if tc.shouldAuth {
				req.AddCookie(&http.Cookie{
					Name:  "auth",
					Value: "secret",
				})
			}

			handler.HandleAuthMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				})).ServeHTTP(rr, req)

			if rr.Code != 200 {
				t.Fatal("expected 200, got", rr.Code)
			}
		})
	}
}

func TestHandler_HandleCreateQuote(t *testing.T) {
	fakeDb := &dbfakes.FakeQuoteStore{}
	tests := map[string]struct {
		wantErr bool
		q       *models.Quote
	}{
		"valid quote": {
			q: models.NewQuote("foo", "bar"),
		},
		"invalid quote": {
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler := handlers.NewHandler(repository.NewQuoteRepository(fakeDb), "secret", true)
			if handler == nil {
				t.Fail()
			}

			req := httptest.NewRequest("POST", "/quotes", nil)
			rr := httptest.NewRecorder()

			if tc.q != nil {
				req.Form = map[string][]string{
					"content": {tc.q.Content},
					"author":  {tc.q.Author},
				}
			}

			if tc.wantErr {
				req.Body = nil
			}

			handler.HandleCreateQuote(rr, req)

			if rr.Code != 200 && !tc.wantErr {
				t.Fatal("expected 200, got", rr.Code)
			} else if rr.Code != 500 && tc.wantErr {
				t.Fatal("expected 500, got", rr.Code)
			}
		})
	}
}

func TestHandler_HandleUpdateQuote(t *testing.T) {
	tests := map[string]struct {
		wantErr bool
		q       *models.Quote
	}{
		"valid quote": {
			q: models.NewQuote("foo", "bar"),
		},
		"invalid quote": {
			wantErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fakeDb := &dbfakes.FakeQuoteStore{}
			fakeDb.QueryByIdReturns(&models.Quote{Id: uuid.New(), Content: "hi", Author: "me"}, nil)
			handler := handlers.NewHandler(repository.NewQuoteRepository(fakeDb), "secret", true)
			if handler == nil {
				t.Fail()
			}

			req := httptest.NewRequest("PUT", "/quotes", nil)
			rr := httptest.NewRecorder()

			if tc.q != nil {
				req.Form = map[string][]string{
					"content": {tc.q.Content},
					"author":  {tc.q.Author},
				}
				req.SetPathValue("id", tc.q.Id.String())
			}

			if tc.wantErr {
				req.Body = nil
			}

			handler.HandleUpdateQuote(rr, req)

			if rr.Code != 200 && !tc.wantErr {
				t.Fatal("expected 200, got", rr.Code)
			} else if rr.Code != 500 && tc.wantErr {
				t.Fatal("expected 500, got", rr.Code)
			}
		})
	}
}
