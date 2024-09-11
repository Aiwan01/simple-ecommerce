package controllers

import "github.com/gofiber/fiber/v2"

func Welcome(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"Message": "Welcome to fiber Api"})
}
