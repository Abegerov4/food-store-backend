package handlers

import (
	"encoding/json"
	"net/http"

	"food-store-backend/repositories"
)

type ProductHandler struct {
	repo *repositories.ProductMongoRepository
}

func NewProductHandler(r *repositories.ProductMongoRepository) *ProductHandler {
	return &ProductHandler{repo: r}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	sort := r.URL.Query().Get("sort")

	products, err := h.repo.GetProducts(category, search, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(products)
}