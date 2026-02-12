package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	"food-store-backend/models"
	"food-store-backend/repositories"
	"food-store-backend/middleware"
	"food-store-backend/utils"
	"go.mongodb.org/mongo-driver/bson"
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
	order.CreatedAt = time.Now()
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
// ADMIN ANALYTICS
func (h *OrderHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	data, err := h.repo.GetAnalytics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}
func (h *OrderHandler) GetRevenueAnalytics(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	data, err := h.repo.GetRevenueByDay()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(data)
}
func (h *OrderHandler) RecalculatePrices(
	productRepo *repositories.ProductMongoRepository,
) (bson.M, error) {

	sales7, err := h.repo.GetSalesLastDays(7)
	if err != nil {
		return nil, err
	}

	sales14, err := h.repo.GetSalesLastDays(14)
	if err != nil {
		return nil, err
	}

	allOrders, _ := h.repo.GetAll()

	updated := 0
	increased := 0
	decreased := 0

	productsMap := make(map[string]int)

	for _, order := range allOrders {
		for _, item := range order.Items {
			productsMap[item.Name] = item.Price
		}
	}

	for name, price := range productsMap {

		newPrice := price
	
		// 📈 Popular product (at least 3 sales in last 7 days)
		if sales7[name] >= 1 {
			newPrice = int(float64(price) * 1.05)
			increased++
		}
	
		// Dead product (low demand)
		if sales14[name] <= 1 {
    	newPrice = int(float64(price) * 0.90)
    	decreased++
		}
	
		// 💘 February seasonal boost
		if time.Now().Month() == time.February {
			newPrice = int(float64(newPrice) * 1.08)
		}
	
		if newPrice != price {
			if newPrice < price {
				productRepo.UpdatePriceWithOriginal(name, price, newPrice)
			} else {
				productRepo.UpdatePriceOnly(name, newPrice)
			}
			updated++
		}
	}

	return bson.M{
		"updated":   updated,
		"increased": increased,
		"decreased": decreased,
	}, nil
}
func (h *OrderHandler) RunDynamicPricing(
	productRepo *repositories.ProductMongoRepository,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		result, err := h.RecalculatePrices(productRepo)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(result)
	}
}