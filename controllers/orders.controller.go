package controllers

import (
	"foodie-service/models"
	"foodie-service/services"
	"foodie-service/types"
	"foodie-service/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrdersController struct {
	services *services.BaseService
	models   *models.BaseModel
}

var ordersController *OrdersController

func NewOrdersController(services *services.BaseService, models *models.BaseModel) *OrdersController {
	if ordersController != nil {
		return ordersController
	}

	return &OrdersController{
		services: services,
		models:   models,
	}
}

func (oc *OrdersController) PlaceOrder(c *fiber.Ctx) error {
	var orderRequest *types.BulkOrdersRequest
	if err := c.BodyParser(&orderRequest); err != nil {
		return utils.ErrorHandler("Error parsing order", err.Error(), fiber.StatusBadRequest, c)
	}

	// Validate the request using our enhanced validator
	if err := utils.Validate(orderRequest); err != nil {
		return utils.ErrorHandler("Invalid order data", err.Error(), fiber.StatusUnprocessableEntity, c)
	}

	userID := c.Locals("userID").(string)

	purchaseDetails, err := oc.services.Orders.PlaceOrder(orderRequest, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorHandler("Product not found", "No product found with the given ID", fiber.StatusNotFound, c)
		}
		return utils.ErrorHandler("Error placing order", err.Error(), fiber.StatusInternalServerError, c)
	}

	return c.JSON(fiber.Map{
		"message": "Order placed successfully",
		"order":   purchaseDetails,
	})
}

func (oc *OrdersController) GetOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	limit := 5
	offset := 0
	var err error
	if c.Get("limit") != "" {
		limit, err = strconv.Atoi(c.Get("limit"))
		if err != nil {
			return utils.ErrorHandler("limit is not integer", "limit must be a valid integer", fiber.StatusBadRequest, c)
		}
	}
	if c.Get("offset") != "" {
		offset, err = strconv.Atoi(c.Get("offset"))
		if err != nil {
			return utils.ErrorHandler("offset is not integer", "offset must be a valid integer", fiber.StatusBadRequest, c)
		}
	}
	purchaseDetails, err := oc.services.Orders.GetPreviousOrders(userID, limit, offset)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return utils.ErrorHandler("no previous orders found", "No order were made previously", fiber.StatusNotFound, c)
		}
		return utils.ErrorHandler("Error placing order", err.Error(), fiber.StatusInternalServerError, c)
	}
	return c.JSON(fiber.Map{
		"message": "Order fetched successfully",
		"order":   purchaseDetails,
	})
}

func (oc *OrdersController) FetchCoupons(c *fiber.Ctx) error {
	coupons, err := oc.services.Coupons.FetchCoupons()
	if err != nil {
		return utils.ErrorHandler("Error fetching coupons", err.Error(), fiber.StatusInternalServerError, c)
	}
	return c.JSON(fiber.Map{
		"message": "Coupons fetched successfully",
		"coupons": coupons,
	})
}