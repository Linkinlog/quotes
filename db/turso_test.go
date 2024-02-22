package db_test

import (
	"database/sql"
	"testing"

	"github.com/Linkinlog/quotes/db"
	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"

	_ "modernc.org/sqlite"
)

func TestTurso_New(t *testing.T) {
	_ = db.NewTursoStore("", "")
}

func TestTurso_QueryById(t *testing.T) {
	tests := map[string]struct {
		quote *models.Quote
	}{
		"valid quote": {
			quote: &models.Quote{
				Id:      uuid.New(),
				Content: "test",
				Author:  "test",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := testDb(t)
			_, err := d.Exec("INSERT INTO quotes (id, content, author, approved) VALUES (?, ?, ?, ?)", tc.quote.Id, tc.quote.Content, tc.quote.Author, tc.quote.Approved())
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			turso := db.Turso{
				Conn: d,
			}
			q, err := turso.QueryById(tc.quote.Id)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if q == nil {
				t.Fatalf("expected quote, got nil")
			}
			if q.Id != tc.quote.Id {
				t.Fatalf("expected id %v, got %v", tc.quote.Id, q.Id)
			}
		})
	}
}

func TestTurso_All(t *testing.T) {
	tests := map[string]struct {
		quotes []*models.Quote
	}{
		"valid quotes": {
			quotes: []*models.Quote{
				{
					Id:      uuid.New(),
					Content: "test",
					Author:  "test",
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := testDb(t)
			turso := db.Turso{
				Conn: d,
			}
			for _, q := range tc.quotes {
				_, err := d.Exec("INSERT INTO quotes (id, content, author, approved) VALUES (?, ?, ?, ?)", q.Id, q.Content, q.Author, q.Approved())
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
			q, err := turso.All()
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if q == nil {
				t.Fatalf("expected quotes, got nil")
			}
			if len(q) != len(tc.quotes) {
				t.Fatalf("expected %d quotes, got %d", len(tc.quotes), len(q))
			}
		})
	}
}

func TestTurso_Insert(t *testing.T) {
	tests := map[string]struct {
		quote *models.Quote
	}{
		"valid quote": {
			quote: &models.Quote{
				Id:      uuid.New(),
				Content: "test",
				Author:  "test",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := testDb(t)
			turso := db.Turso{
				Conn: d,
			}
			err := turso.Insert(tc.quote)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestTurso_Update(t *testing.T) {
	tests := map[string]struct {
		quote *models.Quote
	}{
		"valid quote": {
			quote: &models.Quote{
				Id:      uuid.New(),
				Content: "updated test",
				Author:  "test",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := testDb(t)
			quote := &models.Quote{
				Id:      tc.quote.Id,
				Content: "test",
				Author:  "test",
			}
			_, err := d.Exec("INSERT INTO quotes (id, content, author, approved) VALUES (?, ?, ?, ?)", quote.Id, quote.Content, quote.Author, quote.Approved())
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			turso := db.Turso{
				Conn: d,
			}
			err = turso.Update(tc.quote)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			q, err := turso.QueryById(tc.quote.Id)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if q == nil {
				t.Fatalf("expected quote, got nil")
			}
			if q.Content != tc.quote.Content {
				t.Fatalf("expected content %v, got %v", tc.quote.Content, q.Content)
			}
		})
	}
}

func TestTurso_Delete(t *testing.T) {
	tests := map[string]struct {
		q *models.Quote
	}{
		"valid id": {
			q: &models.Quote{
				Id:      uuid.New(),
				Content: "test",
				Author:  "test",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			d := testDb(t)
			turso := db.Turso{
				Conn: d,
			}
			_, err := d.Exec("INSERT INTO quotes (id, content, author, approved) VALUES (?, ?, ?, ?)", tc.q.Id, tc.q.Content, tc.q.Author, tc.q.Approved())
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			err = turso.Delete(tc.q.Id)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			q, err := turso.QueryById(tc.q.Id)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if q == nil {
				t.Fatalf("expected quote, got nil")
			}
			if q.Approved() {
				t.Fatalf("expected quote to be deleted, got %v", q.Approved())
			}
		})
	}
}

func testDb(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil
	}
	_, err = db.Exec("CREATE TABLE quotes (id TEXT PRIMARY KEY NOT NULL, content TEXT NOT NULL, author TEXT NOT NULL, approved BOOLEAN NOT NULL)")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	t.Cleanup(func() {
		db.Close()
	})
	return db
}
