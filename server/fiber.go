package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())

	api := app.Group("/api", middleware) // /api

	v1 := api.Group("/v1", middleware) // /api/v1

	v1.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Listen(":8246")
}

func middleware(c *fiber.Ctx) error {
	if c.Get("Api_Reader_Secret") == "zxcvasdfqwer1234" {
		return c.Next()
	}
	c.SendStatus(fiber.StatusUnauthorized)
	c.SendString("You are not authorized")
	return nil
}
