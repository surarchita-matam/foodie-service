package utils

import "github.com/gofiber/fiber/v2"

func ErrorHandler(errorType string, errorMessage string, statusCode int, ctx *fiber.Ctx) error {
	ctx.Status(statusCode)
	return ctx.JSON(fiber.Map{
		"errorType": errorType,
		"errorMessage": errorMessage,
		"status": statusCode,
	})
}