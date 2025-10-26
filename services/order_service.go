package services

import (
	"time"

	"food-store-backend/models"
	"food-store-backend/repositories"
)

type OrderService struct {
	orderRepo *repositories.OrderRepository
	foodRepo  *repositories.FoodRepository
}

func NewOrderService(o *repositories.OrderRepository, f *repositories.FoodRepository) *OrderService {
	return &OrderService{o, f}
}

func (s *OrderService) CreateOrder(order models.Order) {
	total := 0.0
	for _, id := range order.FoodIDs {
		if food, ok := s.foodRepo.GetByID(id); ok {
			total += food.Price
		}
	}
	order.Total = total
	order.Status = "processing"

	s.orderRepo.Create(order)

	go s.processOrder(order)
}

func (s *OrderService) processOrder(order models.Order) {
	time.Sleep(2 * time.Second)
	order.Status = "completed"
	s.orderRepo.Update(order)
}

func (s *OrderService) GetOrders() []models.Order {
	return s.orderRepo.GetAll()
}
