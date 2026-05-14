package web

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/artivisi/tech-stack-2026/golang/internal/domain"
	"github.com/artivisi/tech-stack-2026/golang/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Handler struct {
	repo      *repository.Registration
	tmpl      *Templates
	validator *validator.Validate
}

func NewHandler(repo *repository.Registration, tmpl *Templates, v *validator.Validate) *Handler {
	return &Handler{repo: repo, tmpl: tmpl, validator: v}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /", h.ShowForm)
	mux.HandleFunc("POST /register", h.SubmitForm)
	mux.HandleFunc("GET /registrations", h.List)
	mux.HandleFunc("GET /health", h.Health)
}

type formPage struct {
	Values map[string]string
	Errors map[string]string
}

type listPage struct {
	Registrations []domain.Registration
	Count         int
}

func (h *Handler) ShowForm(w http.ResponseWriter, r *http.Request) {
	if err := h.tmpl.Render(w, "form", formPage{
		Values: map[string]string{},
		Errors: map[string]string{},
	}); err != nil {
		log.Printf("render form: %v", err)
	}
}

func (h *Handler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form data", http.StatusBadRequest)
		return
	}

	submitted := map[string]string{
		"email":    r.FormValue("email"),
		"fullName": r.FormValue("fullName"),
		"phone":    r.FormValue("phone"),
	}

	form := RegistrationForm{
		Email:    strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		FullName: strings.TrimSpace(r.FormValue("fullName")),
		Phone:    strings.TrimSpace(r.FormValue("phone")),
	}

	if err := h.validator.Struct(form); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = h.tmpl.Render(w, "form", formPage{
			Values: submitted,
			Errors: CollectErrors(err),
		})
		return
	}

	reg := domain.Registration{
		ID:        uuid.NewString(),
		Email:     form.Email,
		FullName:  form.FullName,
		Phone:     form.Phone,
		CreatedAt: time.Now().UTC(),
	}

	if err := h.repo.Insert(r.Context(), reg); err != nil {
		if errors.Is(err, repository.ErrDuplicateEmail) {
			w.WriteHeader(http.StatusConflict)
			_ = h.tmpl.Render(w, "form", formPage{
				Values: submitted,
				Errors: map[string]string{"email": "email is already registered"},
			})
			return
		}
		log.Printf("insert error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/registrations", http.StatusFound)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	regs, err := h.repo.FindAll(r.Context())
	if err != nil {
		log.Printf("find all error: %v", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if err := h.tmpl.Render(w, "list", listPage{
		Registrations: regs,
		Count:         len(regs),
	}); err != nil {
		log.Printf("render list: %v", err)
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	w.Header().Set("Content-Type", "application/json")
	if err := h.repo.Ping(ctx); err != nil {
		log.Printf("health check failed: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
