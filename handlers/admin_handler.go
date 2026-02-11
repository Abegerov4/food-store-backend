package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"food-store-backend/models"
	"food-store-backend/repositories"
)

type AdminHandler struct {
	repo *repositories.ProductMongoRepository
}

func NewAdminHandler(r *repositories.ProductMongoRepository) *AdminHandler {
	return &AdminHandler{repo: r}
}

// POST
func (h *AdminHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// PUT
func (h *AdminHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/products/")
	var p models.Product

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateByID(id, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DELETE
func (h *AdminHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/admin/products/")

	if err := h.repo.DeleteByID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}