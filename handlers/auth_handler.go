package handlers

import (
	"encoding/json"
	"net/http"

	"food-store-backend/models"
	"food-store-backend/services"
	"food-store-backend/utils"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

// REGISTER
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	user.Role = "user"

	if err := h.service.Register(user); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// LOGIN → JWT
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input models.User

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.Login(input.Email, input.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(
		user.ID.Hex(),
		user.Email,
		user.Role,
	)
	if err != nil {
		http.Error(w, "Token error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}