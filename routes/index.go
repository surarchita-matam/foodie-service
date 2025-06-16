package routes

import (
	"foodie-service/controllers"
	"foodie-service/utils"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	controller := controllers.GetController()

	api := app.Group("/")
	// Public routes
	api.Post("/products", controller.ProductsController.InsertBulkProducts)
	api.Get("/products", controller.ProductsController.GetProducts)
	api.Get("/products/:id", controller.ProductsController.GetProductById)

	// Auth routes
	api.Post("/auth/login", controller.AuthController.Login)
	api.Post("/auth/signup", controller.AuthController.SignUp)

	// Coupons routes
	api.Get("/coupons", controller.OrdersController.FetchCoupons)

	// Protected routes
	secured := api.Group("/orders", utils.ValidateToken())
	secured.Post("/", controller.OrdersController.PlaceOrder)
	secured.Get("/", controller.OrdersController.GetOrders)
}
