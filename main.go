package main

import (
	"context"
	"fmt"
	"go-ecom/routes"
	"log"

	"go-ecom/database"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	port = ":3000"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())
	routes.Router(app)

	client := database.Client

	// Client := database.ConnectWithMongodb()

	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, context.Background())

	if err := app.Listen(port); err != nil {
		log.Fatal("Error in starting server ", err.Error())
	}
	fmt.Println("server running on port :3000")
}
