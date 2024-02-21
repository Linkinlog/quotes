package handlers

import (
	"log/slog"
	"net/http"

	"github.com/Linkinlog/quotes/assets"
	"github.com/Linkinlog/quotes/components"
	"github.com/Linkinlog/quotes/repository"
	"github.com/google/uuid"
)

type Handler struct {
	quoteRepo *repository.QuoteRepository
}

func NewHandler(quoteRepo *repository.QuoteRepository) *Handler {
	return &Handler{quoteRepo: quoteRepo}
}

func (h *Handler) HandleRoutes(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", h.HandleLanding)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Files()))))
	mux.HandleFunc("GET /about", h.HandleAbout)
	mux.HandleFunc("GET /quotes", h.HandleQuotes)
	mux.HandleFunc("GET /quotes/random", h.HandleRandomQuote)
	mux.HandleFunc("GET /quotes/{id}", h.HandleQuoteById)

	return http.ListenAndServe(addr, mux)
}

func (h *Handler) HandleLanding(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.Landing()).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.About()).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleQuotes(w http.ResponseWriter, r *http.Request) {
	q, rErr := h.quoteRepo.All()
	if rErr != nil {
		http.Error(w, "Error getting random quote", http.StatusInternalServerError)
		return
	}
	err := components.Quotes(q).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, rErr := h.quoteRepo.Random()
	if rErr != nil {
		http.Error(w, "Error getting random quote", http.StatusInternalServerError)
		return
	}
	err := components.Quote(q).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleQuoteById(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}
	if id != uuid.Nil {
		q, err := h.quoteRepo.ById(id)
		if err != nil {
			http.Error(w, "Error getting random quote", http.StatusInternalServerError)
			return
		}
		rErr := components.Quote(q).Render(r.Context(), w)
		if rErr != nil {
			slog.Error(rErr.Error())
			return
		}
	}
}
