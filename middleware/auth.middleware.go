package middleware

import (
	"go-ecom/utils"

	"github.com/gofiber/fiber/v2"
)

func RequireAuthValidate(c *fiber.Ctx) error {

	authHeader := c.Get("Authorization")
	token := c.Cookies("jwt")

	if authHeader == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Try to signin first",
		})
	}

	if token == "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Try to signin first",
		})
	}
	id, email, userType, err := utils.VerifyUserToken(token)
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid token",
		})
	}

	c.Locals("id", id)
	c.Locals("email", email)
	c.Locals("userType", userType)

	return c.Next()
}
