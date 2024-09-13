package middleware

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func RegisterCredentialInput(c *fiber.Ctx) error {
	validate := validator.New()

	type UserInputCred struct {
		Username string `json:"username" validate:"required"`
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	var user UserInputCred
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Some field is missing",
		})
	}

	if err := validate.Struct(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid body data",
			"error":   err.Error(),
		})
	}
	return c.Next()
}
