package controllers

import (
	"foodie-service/models"
	"foodie-service/services"
	"foodie-service/types"
	"foodie-service/utils"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	services *services.BaseService
	models   *models.BaseModel
}

var authController *AuthController

func NewAuthController(services *services.BaseService, models *models.BaseModel) *AuthController {
	if authController != nil {
		return authController
	}

	return &AuthController{services: services, models: models}
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var userDetails types.SignInRequest

	if err := c.BodyParser(&userDetails); err != nil {
		return utils.ErrorHandler("Invalid request body", err.Error(), fiber.StatusBadRequest, c)
	}

	signInResponse, err := ac.services.Auth.SignIn(userDetails.Email, userDetails.Password)
	if err != nil {
		return utils.ErrorHandler("Invalid email or password", err.Error(), fiber.StatusUnauthorized, c)
	}
	return c.JSON(fiber.Map{"token": signInResponse.Token})
}

func (ac *AuthController) SignUp(c *fiber.Ctx) error {
	var userDetails types.SignupRequest

	if err := c.BodyParser(&userDetails); err != nil {
		return utils.ErrorHandler("Invalid request body", err.Error(), fiber.StatusBadRequest, c)
	}

	// Validate the request
	if err := utils.Validate(userDetails); err != nil {
		return utils.ErrorHandler("Validation failed", err.Error(), fiber.StatusBadRequest, c)
	}

	signUpResponse, err := ac.services.Auth.SignUp(&userDetails)
	if err != nil {
		return utils.ErrorHandler("Failed to sign up", err.Error(), fiber.StatusInternalServerError, c)
	}
	return c.JSON(fiber.Map{"userId": signUpResponse.UserID})
}
