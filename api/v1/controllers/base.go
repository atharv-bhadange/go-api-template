package controllers

import (
	C "github.com/atharvbhadange/go-api-template/config"
	"github.com/gofiber/fiber/v2"
)

func Health(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"ok":  1,
		"v":   C.Conf.Version,
		"env": C.Conf.Environment,
	})
}
