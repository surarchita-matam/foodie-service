package utils

import (
	"foodie-service/config"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the API key from header
		apiKey := c.Get("api-key")
		if apiKey == "" {
			return ErrorHandler("API key is required", "API key is required", fiber.StatusUnauthorized, c)
		}

		// Remove "Bearer " if present
		tokenString := strings.TrimPrefix(apiKey, "Bearer ")

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Get JWT secret from config
			return []byte(config.GetConfig().JWTSecret), nil
		})

		if err != nil {
			return ErrorHandler("Invalid token", err.Error(), fiber.StatusUnauthorized, c)
		}

		// Type assert the claims
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			return ErrorHandler("Invalid token claims", "Invalid token claims", fiber.StatusUnauthorized, c)
		}

		// Set user ID in locals
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

func GenerateToken(userID string) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.GetConfig().JWTSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
