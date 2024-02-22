package models_test

import (
	"testing"

	"github.com/Linkinlog/quotes/models"
)

func TestNewQuote(t *testing.T) {
	tests := map[string]struct {
		content string
		author  string
	}{
		"valid quote": {
			content: "This is a quote",
			author:  "John Doe",
		},
		"invalid quote": {
			content: "",
			author:  "",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			quote := models.NewQuote(tc.content, tc.author)

			if quote.Approved() {
				t.Fatalf("expected %v; got %v", false, quote.Approved())
			}
			quote.Approve()
			if !quote.Approved() {
				t.Fatalf("expected %v; got %v", false, quote.Approved())
			}

			quote.Disapprove()
			if quote.Approved() {
				t.Fatalf("expected %v; got %v", true, quote.Approved())
			}

			if quote.Content != tc.content {
				t.Errorf("expected %v; got %v", tc.content, quote.Content)
			}

			if quote.Author != tc.author {
				t.Errorf("expected %v; got %v", tc.author, quote.Author)
			}
		})
	}
}
