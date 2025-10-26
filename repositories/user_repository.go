package repositories

import (
	"sync"

	"food-store-backend/models"
)

type UserRepository struct {
	mu    sync.Mutex
	users map[string]models.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{users: make(map[string]models.User)}
}

func (r *UserRepository) Create(user models.User) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
}

func (r *UserRepository) GetAll() []models.User {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := []models.User{}
	for _, u := range r.users {
		result = append(result, u)
	}
	return result
}
