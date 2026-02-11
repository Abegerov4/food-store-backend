package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserEmail string             `bson:"userEmail" json:"userEmail"`
	Items     []Product          `bson:"items" json:"items"`
	Total     int                `bson:"total" json:"total"`
	Status    string             `bson:"status" json:"status"`
}