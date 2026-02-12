package repositories

import (
	"context"

	"food-store-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductMongoRepository struct {
	collection *mongo.Collection
}

func NewProductMongoRepository(db *mongo.Database) *ProductMongoRepository {
	return &ProductMongoRepository{
		collection: db.Collection("products"),
	}
}

// CREATE
func (r *ProductMongoRepository) Create(product models.Product) error {
	_, err := r.collection.InsertOne(context.Background(), product)
	return err
}

// GET
func (r *ProductMongoRepository) GetProducts(
    category, search, sort string,
    page, limit int,
) ([]models.Product, int64, error) {

    filter := bson.M{}

    if category != "" {
        filter["category"] = category
    }

    if search != "" {
        filter["name"] = bson.M{
            "$regex":   search,
            "$options": "i",
        }
    }

    total, err := r.collection.CountDocuments(context.Background(), filter)
    if err != nil {
        return nil, 0, err
    }

    findOptions := options.Find()

    switch sort {
    case "price_asc":
        findOptions.SetSort(bson.D{{Key: "price", Value: 1}})
    case "price_desc":
        findOptions.SetSort(bson.D{{Key: "price", Value: -1}})
    case "name_asc":
        findOptions.SetSort(bson.D{{Key: "name", Value: 1}})
    case "name_desc":
        findOptions.SetSort(bson.D{{Key: "name", Value: -1}})
    }

    skip := (page - 1) * limit
    findOptions.SetSkip(int64(skip))
    findOptions.SetLimit(int64(limit))

    cursor, err := r.collection.Find(context.Background(), filter, findOptions)
    if err != nil {
        return nil, 0, err
    }
    defer cursor.Close(context.Background())

    var products []models.Product
    if err := cursor.All(context.Background(), &products); err != nil {
        return nil, 0, err
    }

    return products, total, nil
}

// UPDATE
func (r *ProductMongoRepository) UpdateByID(id string, product models.Product) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateFields := bson.M{
		"name":     product.Name,
		"price":    product.Price,
		"category": product.Category,
	}

	if product.Image != "" {
		updateFields["image"] = product.Image
	}

	update := bson.M{
		"$set": updateFields,
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		update,
	)

	return err
}

// DELETE
func (r *ProductMongoRepository) DeleteByID(id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.collection.DeleteOne(
		context.Background(),
		bson.M{"_id": objID},
	)

	return err
}
func (r *ProductMongoRepository) UpdatePriceByName(name string, newPrice int) error {

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"name": name},
		bson.M{
			"$set": bson.M{
				"price": newPrice,
			},
		},
	)

	return err
}
func (r *ProductMongoRepository) UpdatePriceWithOriginal(
	name string,
	original int,
	newPrice int,
) error {

	update := bson.M{
		"$set": bson.M{
			"price":         newPrice,
			"originalPrice": original,
		},
	}

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"name": name},
		update,
	)

	return err
}
func (r *ProductMongoRepository) UpdatePriceOnly(
	name string,
	newPrice int,
) error {

	update := bson.M{
		"$set": bson.M{
			"price": newPrice,
		},
		"$unset": bson.M{
			"originalPrice": "",
		},
	}

	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"name": name},
		update,
	)

	return err
}