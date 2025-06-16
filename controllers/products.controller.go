package controllers

import (
	// "foodie-service/dbs"

	"foodie-service/models"
	"foodie-service/services"
	"foodie-service/types"
	"foodie-service/utils"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductsController struct {
	services *services.BaseService
	models   *models.BaseModel
}

var productsController *ProductsController

func NewProductsController(services *services.BaseService, models *models.BaseModel) *ProductsController {
	if productsController != nil {
		return productsController
	}

	return &ProductsController{
		services: services,
		models:   models,
	}
}

func (pc *ProductsController) GetProducts(c *fiber.Ctx) error {

	products, err := pc.models.Products.GetProducts(true)
	if err != nil {
		return utils.ErrorHandler("Error fetching products", err.Error(), fiber.StatusInternalServerError, c)
	}

	return c.JSON(fiber.Map{
		"message":  "Products fetched successfully",
		"products": products,
	})
}

func (pc *ProductsController) GetProductById(c *fiber.Ctx) error {
	id := c.Params("id")

	productId, err := strconv.Atoi(id)
	if err != nil {
		return utils.ErrorHandler("Invalid product ID", "Product ID must be a valid integer", fiber.StatusBadRequest, c)
	}

	product, err := pc.models.Products.GetProductByProductId(strconv.Itoa(productId))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorHandler("Product not found", "No product found with the given ID", fiber.StatusNotFound, c)
		}
		return utils.ErrorHandler("Error fetching product", err.Error(), fiber.StatusInternalServerError, c)
	}
	return c.JSON(fiber.Map{
		"message": "Product fetched successfully",
		"product": product,
	})
}

func (pc *ProductsController) InsertBulkProducts(c *fiber.Ctx) error {

	var bulkProductsRequest types.BulkProductsRequest

	if err := c.BodyParser(&bulkProductsRequest); err != nil {
		return utils.ErrorHandler("Invalid request body", err.Error(), fiber.StatusBadRequest, c)
	}

	if err := utils.Validate(bulkProductsRequest); err != nil {
		return utils.ErrorHandler("Invalid request body", err.Error(), fiber.StatusBadRequest, c)
	}

	err := pc.models.Products.InsertBulkProducts(bulkProductsRequest.Products)
	if err != nil {
		return utils.ErrorHandler("Error inserting products", err.Error(), fiber.StatusInternalServerError, c)
	}

	return c.JSON(fiber.Map{
		"message": "Products inserted successfully",
	})
}
