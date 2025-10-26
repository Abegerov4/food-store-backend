package repositories

import (
	"sync"

	"food-store-backend/models"
)

type FoodRepository struct {
	mu    sync.Mutex
	foods map[string]models.Food
}

func NewFoodRepository() *FoodRepository {
	return &FoodRepository{foods: make(map[string]models.Food)}
}

func (r *FoodRepository) Create(food models.Food) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.foods[food.ID] = food
}

func (r *FoodRepository) GetAll() []models.Food {
	r.mu.Lock()
	defer r.mu.Unlock()

	result := []models.Food{}
	for _, f := range r.foods {
		result = append(result, f)
	}
	return result
}

func (r *FoodRepository) GetByID(id string) (models.Food, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	food, ok := r.foods[id]
	return food, ok
}

func (r *FoodRepository) Update(id string, food models.Food) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.foods[id]; !ok {
		return false
	}
	r.foods[id] = food
	return true
}

func (r *FoodRepository) Delete(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.foods, id)
}
