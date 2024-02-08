package models

import (
	"github.com/google/uuid"
)

type Quote struct {
	id      uuid.UUID
	content string
	author  string
}

func NewQuote(content, author string) *Quote {
	id := uuid.New()
	return &Quote{id: id, content: content, author: author}
}

func (q *Quote) Id() uuid.UUID {
	return q.id
}

func (q *Quote) Content() string {
	return q.content
}

func (q *Quote) Author() string {
	return q.author
}
