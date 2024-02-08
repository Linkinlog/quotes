package handlers

import (
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

func (h *Handler) HandleRoutes() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Files()))))
	http.HandleFunc("/", h.handleLanding)
	http.HandleFunc("/quotes", h.handleQuotes)
	http.HandleFunc("/about", h.handleAbout)
}

func (h *Handler) handleLanding(w http.ResponseWriter, r *http.Request) {
	components.Index(components.Landing()).Render(r.Context(), w)
}

func (h *Handler) handleAbout(w http.ResponseWriter, r *http.Request) {
	components.Index(components.About()).Render(r.Context(), w)
}

func (h *Handler) handleQuotes(w http.ResponseWriter, r *http.Request) {
	sId := r.URL.Query().Get("id")
	if sId == "" {
		q, err := h.quoteRepo.Quotes()
		if err != nil {
			http.Error(w, "Error getting random quote", http.StatusInternalServerError)
			return
		}
		components.Quotes(q).Render(r.Context(), w)
		return
	}
	if sId == "random" {
		q, err := h.quoteRepo.RandomQuote()
		if err != nil {
			http.Error(w, "Error getting random quote", http.StatusInternalServerError)
			return
		}
		components.Quote(q).Render(r.Context(), w)
		return
	}
	id, err := uuid.Parse(sId)
	if err != nil {
		http.Error(w, "Invalid id provided", http.StatusBadRequest)
		return
	}
	if id != uuid.Nil {
		q, err := h.quoteRepo.Quote(id)
		if err != nil {
			http.Error(w, "Error getting random quote", http.StatusInternalServerError)
			return
		}
		components.Quote(q).Render(r.Context(), w)
		return
	}
}
