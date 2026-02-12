package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name          string             `bson:"name" json:"name"`
	Price         int                `bson:"price" json:"price"`
	OriginalPrice int                `bson:"originalPrice,omitempty" json:"originalPrice,omitempty"`
	Image         string             `bson:"image" json:"image"`
	Category      string             `bson:"category" json:"category"`
}