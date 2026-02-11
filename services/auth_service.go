package services

import (

	"food-store-backend/models"
	"food-store-backend/repositories"
)

type AuthService struct {
	repo *repositories.UserRepository
}

func NewAuthService(r *repositories.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Register(user models.User) error {
	return s.repo.Create(user)
}

func (s *AuthService) Login(email, password string) (models.User, error) {
	return s.repo.FindByCredentials(email, password)
}