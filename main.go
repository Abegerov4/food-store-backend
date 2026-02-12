package main

import (
	"log"
	"net/http"

	"food-store-backend/config"
	"food-store-backend/handlers"
	"food-store-backend/middleware"
	"food-store-backend/repositories"
	"food-store-backend/services"
)

func main() {

	//  Mongo
	config.ConnectMongo()

	//  Repositories
	productRepo := repositories.NewProductMongoRepository(config.DB)
	userRepo := repositories.NewUserRepository()
	orderRepo := repositories.NewOrderMongoRepository(config.DB) 

	//  Services
	authService := services.NewAuthService(userRepo)

	//  Handlers
	productHandler := handlers.NewProductHandler(productRepo)
	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(productRepo)
	orderHandler := handlers.NewOrderHandler(orderRepo) 

	//  PUBLIC ROUTES

	http.HandleFunc("/api/products", productHandler.GetProducts)

	http.HandleFunc("/register", authHandler.Register)
	http.HandleFunc("/login", authHandler.Login)
	analyticsRouter := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			http.HandlerFunc(orderHandler.GetAnalytics),
		),
	)
	
	http.Handle("/api/admin/analytics", analyticsRouter)
	revenueRouter := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			http.HandlerFunc(orderHandler.GetRevenueAnalytics),
		),
	)
	
	http.Handle("/api/admin/analytics/revenue", revenueRouter)
	//  ORDERS (checkout + my orders)


	ordersRouter := middleware.AuthMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				orderHandler.CreateOrder(w, r)
			case http.MethodGet:
				orderHandler.GetMyOrders(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		}),
	)
	
	http.Handle("/api/orders", ordersRouter)

	//  ADMIN ROUTES

	adminRouter := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				switch r.Method {
				case http.MethodPost:
					adminHandler.CreateProduct(w, r)
				case http.MethodPut:
					adminHandler.UpdateProduct(w, r)
				case http.MethodDelete:
					adminHandler.DeleteProduct(w, r)
				default:
					http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
				}
			}),
		),
	)
	// ADMIN UPDATE ORDER STATUS
	adminOrdersRouter := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			http.HandlerFunc(orderHandler.UpdateOrderStatus),
		),
	)
	// ADMIN GET ALL ORDERS
	adminGetOrders := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			http.HandlerFunc(orderHandler.GetAllOrders),
		),
	)
	pricingRouter := middleware.AuthMiddleware(
		middleware.AdminMiddleware(
			orderHandler.RunDynamicPricing(productRepo),
		),
	)
	
	http.Handle("/api/admin/pricing/recalculate", pricingRouter)

http.Handle("/api/admin/orders/all", adminGetOrders)

http.Handle("/api/admin/orders", adminOrdersRouter)

	http.Handle("/api/admin/products", adminRouter)
	http.Handle("/api/admin/products/", adminRouter)

	// STATIC FILES

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	log.Println(" Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}