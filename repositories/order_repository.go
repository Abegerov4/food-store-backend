package repositories

import (
	"sync"

	"food-store-backend/models"
)

type OrderRepository struct {
	mu     sync.Mutex
	orders map[string]models.Order
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{orders: make(map[string]models.Order)}
}

func (r *OrderRepository) Create(order models.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
}

func (r *OrderRepository) Update(order models.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
}

func (r *OrderRepository) GetAll() []models.Order {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []models.Order{}
	for _, o := range r.orders {
		result = append(result, o)
	}
	return result
}
