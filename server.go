package main

import (
	"context"
	"fmt"
	"foodie-service/controllers"

	// "foodie-service/coupons"
	"foodie-service/database"
	"foodie-service/models"
	"foodie-service/routes"
	"foodie-service/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

// CreateServer initializes and starts the Fiber server on port 3000
func CreateServer(ctx context.Context) {

	app := fiber.New(fiber.Config{
		AppName: "Foodie Service v1.0.0",
	})

	// Add logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${method} | ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	app.Use(pprof.New())

	// Add panic recovery middleware
	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("Recovered from panic: %v\n", r)
				c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error",
				})
			}
		}()
		return c.Next()
	})

	mongoClientPrimary, err := database.MongoClient("primary")
	if err != nil {
		fmt.Println("Connection to primary mongo instance could not be established", err)
	}
	mongoClientSecondary, err := database.MongoClient("secondary")
	if err != nil {
		fmt.Println("Connection to secondary mongo instance could not be established", err)
	}

	models := models.NewBaseModel(mongoClientPrimary, mongoClientSecondary)
	services := services.NewBaseService(models)
	controllers.NewBaseController(services, models)

	// Initialize coupon package
	fmt.Println("Loading coupons into memory...")
	if err := services.Coupons.Init(ctx); err != nil {
		fmt.Printf("Failed to load coupons: %v\n", err)
		return
	}
	fmt.Println("Coupons loaded successfully.")
	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	routes.SetupRoutes(app)

	// Run the server in a goroutine so it doesn't block
	go func() {
		if err := app.Listen(":3000"); err != nil {
			fmt.Printf("Server stopped: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	fmt.Println("server stopped")
	if err := app.Shutdown(); err != nil {
		fmt.Println("server Shutdown Failed:%+v", err)
	}
	fmt.Println("server exited properly")
}
