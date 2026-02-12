package handlers

import (
	"encoding/json"
	"net/http"

	"food-store-backend/models"
	"food-store-backend/repositories"
	"food-store-backend/middleware"
	"food-store-backend/utils"
)

type OrderHandler struct {
	repo *repositories.OrderMongoRepository
}

func NewOrderHandler(r *repositories.OrderMongoRepository) *OrderHandler {
	return &OrderHandler{repo: r}
}

//  CREATE ORDER (checkout)
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var order models.Order

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//  получаем email из JWT
	claimsValue := r.Context().Value(middleware.UserContextKey)
	if claimsValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userClaims := claimsValue.(*utils.Claims)

	order.UserEmail = userClaims.Email
	order.Status = "processing"

	err = h.repo.Create(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(order)
}

//  GET ORDERS BY USER EMAIL
func (h *OrderHandler) GetMyOrders(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Content-Type", "application/json")

    claims := r.Context().Value(middleware.UserContextKey)
    if claims == nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    userClaims := claims.(*utils.Claims)
    email := userClaims.Email

    orders, err := h.repo.GetByEmail(email)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(orders)
}
//  ADMIN: UPDATE ORDER STATUS
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Order ID required", http.StatusBadRequest)
		return
	}

	var body struct {
		Status string `json:"status"`
	}

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.repo.UpdateStatus(id, body.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
//  ADMIN: GET ALL ORDERS
func (h *OrderHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	orders, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(orders)
}