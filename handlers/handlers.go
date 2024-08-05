package handlers

import (
	"fmt"
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
	mux.HandleFunc("GET /robots.txt", h.HandleAuthMiddleware(h.HandleRobots))
	mux.HandleFunc("GET /about", h.HandleAuthMiddleware(h.HandleAbout))
	mux.HandleFunc("GET /add", h.HandleAuthMiddleware(h.HandleCreateForm))
	mux.HandleFunc("GET /quotes", h.HandleAuthMiddleware(h.HandleQuotes))
	mux.HandleFunc("GET /quotes/random", h.HandleAuthMiddleware(h.HandleRandomQuote))
	mux.HandleFunc("GET /quotes/{id}", h.HandleAuthMiddleware(h.HandleQuoteById))
	mux.HandleFunc("GET /quotes/{id}/edit", h.HandleAuthMiddleware(h.HandleEditForm))
	mux.HandleFunc("POST /quotes/{id}/approve", h.HandleAuthMiddleware(h.HandleApproveQuote))
	mux.HandleFunc("DELETE /quotes/{id}", h.HandleAuthMiddleware(h.HandleDeleteQuote))
	mux.HandleFunc("POST /quotes", h.HandleAuthMiddleware(h.HandleCreateQuote))
	mux.HandleFunc("PUT /quotes/{id}", h.HandleAuthMiddleware(h.HandleUpdateQuote))

	fmt.Printf("Listening on %s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func (h *Handler) HandleLanding(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.Landing(), "WisePup", "Quotes from a wise pup").Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleRobots(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("User-agent: *\nDisallow: /"))
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleAbout(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.About(), "WisePup | About", "Made with HTMX, Go, Love, and more...").Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleCreateForm(w http.ResponseWriter, r *http.Request) {
	err := components.Index(components.AddQuote(false), "WisePup | Create", "All quotes are appreciated").Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleQuotes(w http.ResponseWriter, r *http.Request) {
	q, rErr := h.quoteRepo.All()
	if rErr != nil {
		http.Error(w, "Error getting random quote", http.StatusInternalServerError)
		slog.Error(rErr.Error())
		return
	}
	err := components.Index(components.Quotes(q, h.admin), "WisePup | All Quotes", "Every quote, so far").Render(r.Context(), w)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func (h *Handler) HandleRandomQuote(w http.ResponseWriter, r *http.Request) {
	q, rErr := h.quoteRepo.Random()
	if rErr != nil {
		http.Error(w, "Error getting random quote", http.StatusInternalServerError)
		slog.Error(rErr.Error())
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
			slog.Error(err.Error())
			return
		}
		rErr := components.Index(components.Quote(q, h.admin), "WisePup", fmt.Sprintf("\"%s\"\n- %s", q.Content, q.Author)).Render(r.Context(), w)
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
		slog.Error(err.Error())
		return
	}
	content := r.FormValue("content")
	author := r.FormValue("author")

	err = h.quoteRepo.Insert(content, author)
	if err != nil {
		http.Error(w, "Error inserting quote", http.StatusInternalServerError)
		slog.Error(err.Error())
		return
	}
	rErr := components.Index(components.AddQuote(true), "WisePup", fmt.Sprintf("\"%s\"\n- %s", content, author)).Render(r.Context(), w)
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
		slog.Error(err.Error())
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
		slog.Error(err.Error())
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
		slog.Error(err.Error())
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	aErr := h.quoteRepo.Approve(id)
	if aErr != nil {
		http.Error(w, "Error approving quote", http.StatusInternalServerError)
		slog.Error(aErr.Error())
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
		slog.Error(err.Error())
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}

	dErr := h.quoteRepo.Delete(id)
	if dErr != nil {
		http.Error(w, "Error approving quote", http.StatusInternalServerError)
		slog.Error(dErr.Error())
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
		slog.Error(err.Error())
		return
	}
	if q == nil {
		http.Error(w, "Quote not found", http.StatusNotFound)
		return
	}
	err = components.Index(components.EditQuote(q, false), "WisePup - Edit", "Edit your quote!").Render(r.Context(), w)
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
			slog.Error(err.Error())
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
