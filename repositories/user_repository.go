package repositories

import (
	"context"
	"errors"

	"food-store-backend/config"
	"food-store-backend/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// CREATE USER
func (r *UserRepository) Create(user models.User) error {
	collection := config.DB.Collection("users")

	// check if exists
	count, _ := collection.CountDocuments(
		context.Background(),
		bson.M{"email": user.Email},
	)

	if count > 0 {
		return errors.New("user already exists")
	}

	_, err := collection.InsertOne(context.Background(), user)
	return err
}

// FIND BY EMAIL + PASSWORD
func (r *UserRepository) FindByCredentials(email, password string) (models.User, error) {
	collection := config.DB.Collection("users")

	var user models.User
	err := collection.FindOne(
		context.Background(),
		bson.M{
			"email":    email,
			"password": password,
		},
	).Decode(&user)

	if err != nil {
		return user, errors.New("invalid credentials")
	}

	return user, nil
}