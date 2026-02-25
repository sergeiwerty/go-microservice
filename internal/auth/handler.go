package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	ctx         context.Context
	authService AuthService
}

func NewHandler(ctx context.Context, as AuthService) *Handler {
	return &Handler{
		authService: as,
		ctx:         ctx,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	user, err := h.authService.Register(r.Context(), req.Username, req.Password)

	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			http.Error(w, "User already exists", http.StatusConflict)

			return
		}
		http.Error(w, "Server error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)

		return
	}
	defer r.Body.Close()

	success, _ := h.authService.Login(r.Context(), req.Username, req.Password)

	if !success {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful!",
	})
}
