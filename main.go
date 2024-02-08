package main

import (
	"log/slog"
	"net/http"

	"github.com/Linkinlog/quotes/db"
	"github.com/Linkinlog/quotes/handlers"
	"github.com/Linkinlog/quotes/repository"
)

const addr = ":8080"

func main() {
	store := db.NewInMemoryStore(repository.MakeQuotes())
	repo := repository.NewQuoteRepository(store)

	handlers.NewHandler(repo).HandleRoutes()

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("ListenAndServe:", err)
	}
	return
}
