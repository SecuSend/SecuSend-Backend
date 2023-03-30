package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func PostNote() fiber.Handler {
	return func(c *fiber.Ctx) error {
		test := c.Query("test")

		c.Status(fiber.StatusOK)
		return c.JSON(fiber.Map{"test": test})
	}
}
