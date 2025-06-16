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
	api.Get("/products", controller.ProductsController.GetProducts)
	api.Get("/products/:id", controller.ProductsController.GetProductById)
	// api.Post("/product", controllers.CreateProduct)

	api.Post("/login", controller.AuthController.Login)
	api.Post("/signup", controller.AuthController.SignUp)

	api.Post("/load-products", controller.ProductsController.InsertBulkProducts)

	api.Get("/fetch-coupons", controller.OrdersController.FetchCoupons)

	// Protected routes
	secured := api.Group("/orders", utils.ValidateToken())
	secured.Post("/", controller.OrdersController.PlaceOrder)
	secured.Get("/", controller.OrdersController.GetOrders)
}
