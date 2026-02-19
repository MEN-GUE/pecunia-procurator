package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	ID          int     `json:"_id, omitempty" bson:"_id",emitempty`
	Type        string  `json:"type"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, world")

	// Fetch .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("File .env not found", err)
	}

	// Read .env file to connecto to mongodb
	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB Atlas")

	collection = client.Database("go_db").Collection("transactions")

	// Start application
	app := fiber.New()

	app.Get("/api/trans", getTransactions)
	app.Post("/api/trans", postTransaction)

	// Configure port
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getTransactions(c fiber.Ctx) error {
	var transactions []Transaction

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var transaction Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return err
		}

		transactions = append(transactions, transaction)
	}

	return c.JSON(transactions)
}

func postTransaction(c fiber.Ctx) error {
	transaction := new(Transaction)

	if err := c.Bind().Body(transaction); err != nil {
		return err
	}

	if transaction.Amount == 0.0 {
		return c.Status(400).JSON(fiber.Map{"error": "transaction amount must be grater than 0"})
	}

	insertResult, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return err
	}

	transaction.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(200).JSON(transaction)

}
