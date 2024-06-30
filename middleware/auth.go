package middleware

import (
	"synapsis/models"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateUser(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	tokenString = tokenString[7:] // Remove "Bearer " prefix

	if tokenString == "" {
		return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Missing or malformed JWT")
	}

	claims, err := utils.ValidateJWT(tokenString)
	if err != nil {

		return utils.ErrorResponseJSON(c, fiber.StatusUnauthorized, "Invalid or expired JWT")
	}

	c.Locals("userID", claims.UserID)
	c.Locals("role", claims.Role)
	return c.Next()
}

func AdminOnly(c *fiber.Ctx) error {
	if c.Locals("role") != models.RoleAdmin {
		return utils.ErrorResponseJSON(c, fiber.StatusForbidden, "You are not authorized to access this resource")
	}
	return c.Next()
}
