package repositories

import (
	"context"
	"time"
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
// ADMIN ANALYTICS
func (r *OrderMongoRepository) GetAnalytics() (bson.M, error) {

	ctx := context.Background()

	// 1️⃣ Total revenue (completed only)
	revenuePipeline := mongo.Pipeline{
		{{"$match", bson.D{{"status", "completed"}}}},
		{{"$group", bson.D{
			{"_id", nil},
			{"totalRevenue", bson.D{{"$sum", "$total"}}},
			{"totalOrders", bson.D{{"$sum", 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, revenuePipeline)
	if err != nil {
		return nil, err
	}

	var revenueResult []bson.M
	if err := cursor.All(ctx, &revenueResult); err != nil {
		return nil, err
	}

	totalRevenue := 0
	totalOrders := 0

	if len(revenueResult) > 0 {
		totalRevenue = int(revenueResult[0]["totalRevenue"].(int32))
		totalOrders = int(revenueResult[0]["totalOrders"].(int32))
	}

	// 2️⃣ Top products
	topProductsPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"status", "completed"}}}},
		{{"$unwind", "$items"}},
		{{"$group", bson.D{
			{"_id", "$items.name"},
			{"quantity", bson.D{{"$sum", 1}}},
		}}},
		{{"$sort", bson.D{{"quantity", -1}}}},
		{{"$limit", 5}},
	}

	cursor2, err := r.collection.Aggregate(ctx, topProductsPipeline)
	if err != nil {
		return nil, err
	}

	var topProducts []bson.M
	if err := cursor2.All(ctx, &topProducts); err != nil {
		return nil, err
	}

	return bson.M{
		"totalRevenue": totalRevenue,
		"totalOrders":  totalOrders,
		"topProducts":  topProducts,
	}, nil
}
func (r *OrderMongoRepository) GetRevenueByDay() ([]bson.M, error) {

	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"status", "completed"}}}},
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m-%d"},
					{"date", "$createdAt"},
				}},
			}},
			{"revenue", bson.D{{"$sum", "$total"}}},
		}}},
		{{"$sort", bson.D{{"_id", 1}}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}
func (r *OrderMongoRepository) GetSalesLastDays(days int) (map[string]int, error) {

	ctx := context.Background()

	pipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"status", "completed"},
			{"createdAt", bson.D{
				{"$gte", primitive.NewDateTimeFromTime(
					time.Now().AddDate(0, 0, -days),
				)},
			}},
		}}},
		{{"$unwind", "$items"}},
		{{"$group", bson.D{
			{"_id", "$items.name"},
			{"count", bson.D{{"$sum", 1}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}

	var results []bson.M
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	sales := make(map[string]int)

	for _, r := range results {
		name := r["_id"].(string)
		count := int(r["count"].(int32))
		sales[name] = count
	}

	return sales, nil
}