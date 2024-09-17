package controllers

import (
	"context"
	"go-ecom/database"
	"go-ecom/models"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var productCollection = database.OpenCollection(database.Client, "products")

func CreateNewProducts(c *fiber.Ctx) error {
	gofakeit.Seed(0)
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	userType, err := c.Locals("userType").(string)

	if userType != "ADMIN" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Only Admin can add product",
			"data":    err,
		})
	}

	var product models.Product
	product.ID = primitive.NewObjectID()
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()
	product.Name = gofakeit.Name()
	product.Description = gofakeit.Sentence(10)
	product.Price = gofakeit.Price(10, 1000)
	product.Images = []string{gofakeit.ImageURL(100, 100)}
	product.AvailableQuantity = gofakeit.Number(1, 100)

	filter := bson.M{"name": product.ID}
	if existingProduct, err := productCollection.FindOne(ctx, filter).DecodeBytes(); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Product already exist",
			"data":    existingProduct,
		})
	}

	if _, err := productCollection.InsertOne(ctx, product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Product already exist",
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Product created successfully",
		"data":    product,
	})
}

func GetProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	productId := c.Params("id")
	product_id, err := primitive.ObjectIDFromHex(productId)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Unable to get product id",
			"data":    err,
		})
	}

	var product models.Product
	if err := productCollection.FindOne(ctx, bson.M{"_id": product_id}).Decode(&product); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Unable to fetch product ",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "Success",
		"message": "Successfully fetched product.",
		"data":    product,
	})

}

func GetAllProducts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	productRowList, err := productCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Getting error in product fetch",
			"data":    err,
		})
	}

	var products []bson.M
	if err := productRowList.All(ctx, &products); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed parse data",
			"data":    err,
		})
	}
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "products list fetched successfully",
		"data":    products,
	})
}

func UpdateProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	var productData bson.M
	if err := c.BodyParser(&productData); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Fail to parse body",
			"data":    err,
		})
	}

	userType := c.Locals("userType").(string)
	if userType != "ADMIN" {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin can update product",
			"data":    nil,
		})
	}

	id := c.Params("id")
	productId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "cant convert id into object id",
			"data":    err,
		})
	}

	filter := bson.M{"_id": productId}
	product, err := productCollection.FindOneAndUpdate(ctx, filter, bson.M{"$set": productData}).DecodeBytes()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Fail to update product",
			"data":    err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "successfully updated product",
		"data":    product.String(),
	})
}

func DeleteProduct(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	id := c.Params("id")
	productId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Can't convert id into primit data",
			"data":    err,
		})
	}

	userType := c.Locals("userType")
	if userType == "" {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user data modeling",
			"data":    nil,
		})
	}

	if userType != "ADMIN" {
		return c.Status(403).JSON(fiber.Map{
			"status":  "error",
			"message": "Only admin can delete product",
			"data":    nil,
		})
	}

	filter := bson.M{"_id": productId}
	singleResult := productCollection.FindOneAndDelete(ctx, filter)
	var product models.Product
	if err := singleResult.Decode(&product); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Something when wrong on delete",
			"data":    err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "Product deleted successful",
		"data":    nil,
	})
}
