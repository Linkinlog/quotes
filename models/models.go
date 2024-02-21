package models

import (
	"github.com/google/uuid"
)

type Quote struct {
	Id      uuid.UUID
	Content string
	Author  string
}

func NewQuote(content, author string) *Quote {
	id := uuid.New()
	return &Quote{Id: id, Content: content, Author: author}
}
