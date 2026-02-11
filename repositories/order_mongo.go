package repositories

import (
	"context"

	"food-store-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderMongoRepository struct {
	collection *mongo.Collection
}

func NewOrderMongoRepository(db *mongo.Database) *OrderMongoRepository {
	return &OrderMongoRepository{
		collection: db.Collection("orders"),
	}
}

// CREATE ORDER
func (r *OrderMongoRepository) Create(order models.Order) error {
	_, err := r.collection.InsertOne(context.Background(), order)
	return err
}

func (r *OrderMongoRepository) GetByEmail(email string) ([]models.Order, error) {

	filter := bson.M{"userEmail": email}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
// UPDATE ORDER STATUS
func (r *OrderMongoRepository) UpdateStatus(id string, status string) error {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	return err
}
// GET ALL ORDERS (ADMIN)
func (r *OrderMongoRepository) GetAll() ([]models.Order, error) {

	cursor, err := r.collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var orders []models.Order
	if err := cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}

	return orders, nil
}