package controllers

import (
	"foodie-service/models"
	"foodie-service/services"
	"foodie-service/types"
	"foodie-service/utils"

	"github.com/gofiber/fiber/v2"
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

	userID := c.Locals("userID").(string)

	purchaseDetails, err := oc.services.Orders.PlaceOrder(orderRequest, userID)
	if err != nil {
		return utils.ErrorHandler("Error placing order", err.Error(), fiber.StatusInternalServerError, c)
	}

	return c.JSON(fiber.Map{
		"message": "Order placed successfully",
		"order":   purchaseDetails,
	})
}
