package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/Linkinlog/quotes/models"
	"github.com/google/uuid"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Turso struct {
	Token    string
	Database string
	Conn     *sql.DB
}

func NewTursoStore(token, database string) QuoteStore {
	dataSource := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", database, token)
	db, err := sql.Open("libsql", dataSource)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to open db %s: %s", dataSource, err))
		os.Exit(1)
	}
	return &Turso{Token: token, Database: database, Conn: db}
}

func (t *Turso) Insert(q *models.Quote) error {
	_, err := t.Conn.Exec("INSERT INTO quotes (id, content, author, approved) VALUES (?, ?, ?, ?)", q.Id, q.Content, q.Author, q.Approved())
	if err != nil {
		return err
	}
	return nil
}

func (t *Turso) QueryById(id uuid.UUID) (*models.Quote, error) {
	var q models.Quote
	var approved bool

	r := t.Conn.QueryRow("SELECT * FROM quotes WHERE id = ?", id)

	if err := r.Scan(&q.Id, &q.Content, &q.Author, &approved); err != nil {
		return nil, err
	}
	if approved && !q.Approved() {
		q.Approve()
	}

	return &q, nil
}

func (t *Turso) All() ([]*models.Quote, error) {
	r, err := t.Conn.Query("SELECT * FROM quotes")
	if err != nil {
		return nil, err
	}

	quotes := make([]*models.Quote, 0)
	for r.Next() {
		var q models.Quote
		var approved bool
		if err := r.Scan(&q.Id, &q.Content, &q.Author, &approved); err != nil {
			return nil, err
		}
		if approved && !q.Approved() {
			q.Approve()
		}
		quotes = append(quotes, &q)
	}
	return quotes, nil
}

func (t *Turso) Update(q *models.Quote) error {
	_, err := t.Conn.Exec("UPDATE quotes SET content = ?, author = ?, approved = ? WHERE id = ?", q.Content, q.Author, q.Approved(), q.Id)
	if err != nil {
		return err
	}
	return nil
}

func (t *Turso) Delete(id uuid.UUID) error {
	_, err := t.Conn.Exec("UPDATE quotes SET approved = ? WHERE id = ?", false, id)
	if err != nil {
		return err
	}
	return nil
}
