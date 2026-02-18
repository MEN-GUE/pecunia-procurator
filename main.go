package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Transaction struct {
	ID          int     `json:"id"`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

func main() {
	// Welcome message to debug test
	fmt.Println("Hello world! This is Pecunia-Procurator")

	// Create the web application
	app := fiber.New()
	// Load and Read the .env for getting port
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}

	PORT := os.Getenv("PORT")

	// Save transaction in local memmory
	transactions := []Transaction{}

	// Get all transactions
	app.Get("/api/trans", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(transactions)
	})

	// Post a transaction
	app.Post("/api/trans", func(c *fiber.Ctx) error {
		transaction := &Transaction{}

		if err := c.BodyParser(transaction); err != nil {
			return err
		}

		if transaction.Amount == 0.0 {
			return c.Status(400).JSON(fiber.Map{"error": "Enter at least a valid amount"})
		}

		transaction.ID = len(transactions) + 1
		transactions = append(transactions, *transaction)

		return c.Status(201).JSON(transaction)
	})

	// TODO: Delete a transaction
	//app.Delete("/api/trans/:id", func(c *fiber.Ctx) error {
	//	id := c.Params("id")
	//})

	// Listen to server port
	log.Fatal(app.Listen(":" + PORT))
}
