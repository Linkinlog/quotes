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
	secret    string
	admin     bool
}

func NewHandler(quoteRepo *repository.QuoteRepository, secret string, admin bool) *Handler {
	return &Handler{quoteRepo: quoteRepo, secret: secret, admin: admin}
}

func (h *Handler) HandleRoutes(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(assets.Files()))))

	mux.HandleFunc("GET /", h.HandleAuthMiddleware(h.HandleLanding))
	mux.HandleFunc("GET /about", h.HandleAuthMiddleware(h.HandleAbout))
	mux.HandleFunc("GET /add", h.HandleAuthMiddleware(h.HandleCreateForm))
	mux.HandleFunc("GET /quotes", h.HandleAuthMiddleware(h.HandleQuotes))
	mux.HandleFunc("GET /quotes/random", h.HandleAuthMiddleware(h.HandleRandomQuote))
	mux.HandleFunc("GET /quotes/{id}", h.HandleAuthMiddleware(h.HandleQuoteById))
	mux.HandleFunc("GET /quotes/{id}/edit", h.HandleAuthMiddleware(h.HandleEditForm))
	mux.HandleFunc("POST /quotes/{id}/toggle", h.HandleAuthMiddleware(h.HandleApproveQuote))
	mux.HandleFunc("DELETE /quotes/{id}", h.HandleAuthMiddleware(h.HandleDeleteQuote))
	mux.HandleFunc("POST /quotes", h.HandleAuthMiddleware(h.HandleCreateQuote))
	mux.HandleFunc("PUT /quotes/{id}", h.HandleAuthMiddleware(h.HandleUpdateQuote))

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

func (h *Handler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.AddQuote(false)).Render(r.Context(), w)
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
	err := components.Index(components.Quotes(q, h.admin)).Render(r.Context(), w)
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
	err := components.Quote(q, h.admin).Render(r.Context(), w)
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
		rErr := components.Index(components.Quote(q, h.admin)).Render(r.Context(), w)
		if rErr != nil {
			slog.Error(rErr.Error())
			return
		}
	}
}

func (h *Handler) HandleCreateQuote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}
	err = h.quoteRepo.Insert(r.FormValue("content"), r.FormValue("author"))
	if err != nil {
		http.Error(w, "Error inserting quote", http.StatusInternalServerError)
		return
	}
	rErr := components.Index(components.AddQuote(true)).Render(r.Context(), w)
	if rErr != nil {
		slog.Error(rErr.Error())
		return
	}
}

func (h *Handler) HandleUpdateQuote(w http.ResponseWriter, r *http.Request) {
	if !h.admin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}
	q, err := h.quoteRepo.Update(id, r.FormValue("content"), r.FormValue("author"))
	if err != nil {
		http.Error(w, "Error updating quote", http.StatusInternalServerError)
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}
	err = components.Quote(q, h.admin).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleApproveQuote(w http.ResponseWriter, r *http.Request) {
	if !h.admin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}
	q, err := h.quoteRepo.ById(id)
	if err != nil {
		http.Error(w, "Error getting quote", http.StatusInternalServerError)
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	q.Approve()
	err = components.Quote(q, h.admin).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleDeleteQuote(w http.ResponseWriter, r *http.Request) {
	if !h.admin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}
	q, err := h.quoteRepo.ById(id)
	if err != nil {
		http.Error(w, "Error getting quote", http.StatusInternalServerError)
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	q.Disapprove()
	err = components.Quote(q, h.admin).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleEditForm(w http.ResponseWriter, r *http.Request) {
	if !h.admin {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil || id == uuid.Nil {
		http.Error(w, "invalid Id", http.StatusBadRequest)
		return
	}
	q, err := h.quoteRepo.ById(id)
	if err != nil {
		http.Error(w, "Error getting quote", http.StatusInternalServerError)
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}
	err = components.Index(components.EditQuote(q, false)).Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.admin = false
		val, err := r.Cookie("auth")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		if val != nil && val.Value == h.secret && h.secret != "" {
			h.admin = true
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
