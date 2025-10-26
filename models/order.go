package models

type Order struct {
	ID      string   `json:"id"`
	FoodIDs []string `json:"food_ids"`
	Total   float64  `json:"total"`
	Status  string   `json:"status"`
}
