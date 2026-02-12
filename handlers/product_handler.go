package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"food-store-backend/repositories"
)

type ProductHandler struct {
	repo *repositories.ProductMongoRepository
}

func NewProductHandler(r *repositories.ProductMongoRepository) *ProductHandler {
	return &ProductHandler{repo: r}
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {

    category := r.URL.Query().Get("category")
    search := r.URL.Query().Get("search")
    sort := r.URL.Query().Get("sort")

    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page := 1
    limit := 12

    if pageStr != "" {
        p, _ := strconv.Atoi(pageStr)
        if p > 0 {
            page = p
        }
    }

    if limitStr != "" {
        l, _ := strconv.Atoi(limitStr)
        if l > 0 {
            limit = l
        }
    }

    products, total, err := h.repo.GetProducts(category, search, sort, page, limit)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "data":  products,
        "total": total,
    }

    json.NewEncoder(w).Encode(response)
}