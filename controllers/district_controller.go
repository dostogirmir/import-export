package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func ImportDistricts(c *fiber.Ctx) error {
	return c.SendString("ImportDistricts............")
}