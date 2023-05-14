package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Product struct {
	gorm.Model
	Name        string
	Description string
	Price       float64
}

func main() {
	app := fiber.New()

	db, err := gorm.Open("sqlite3", "trying.sqlite")
	if err != nil {
		log.Printf("failed to connect database: %v", err)
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&Product{})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/products", func(c *fiber.Ctx) error {
		var products []Product
		db.Find(&products)
		return c.JSON(products)
	})

	app.Post("/products", func(c *fiber.Ctx) error {
		product := new(Product)
		if err := c.BodyParser(product); err != nil {
			return err
		}
		db.Create(&product)
		return c.JSON(product)
	})

	app.Listen(":3000")
}
