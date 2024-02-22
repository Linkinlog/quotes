package models

import (
	"github.com/google/uuid"
)

type Quote struct {
	Id       uuid.UUID
	Content  string
	Author   string
	approved bool
}

func NewQuote(content, author string) *Quote {
	id := uuid.New()
	return &Quote{Id: id, Content: content, Author: author}
}

func (q *Quote) Approved() bool {
	return q.approved
}

func (q *Quote) Approve() {
	q.approved = true
}

func (q *Quote) Disapprove() {
	q.approved = false
}
