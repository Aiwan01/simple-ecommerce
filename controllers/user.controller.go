package controllers

import (
	"context"
	"go-ecom/database"
	"go-ecom/models"
	"go-ecom/utils"
	"os"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenConnection(database.Client, "users")

func Signup(c *fiber.Ctx) error {

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*5)
	defer cancel()

	var user models.User
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	var userAddress models.Address
	userAddress.ZipCode = gofakeit.Zip()
	userAddress.City = gofakeit.City()
	userAddress.State = gofakeit.State()
	userAddress.Country = gofakeit.Country()
	userAddress.Street = gofakeit.Street()
	userAddress.HouseNumber = gofakeit.StreetNumber()
	user.Address = userAddress
	user.Orders = make([]models.Order, 0)
	user.UserCart = make([]models.ProductsToOrder, 0)

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "Invalid user data for model binding",
			"data":    err.Error(),
		})
	}

	if user.Password == os.Getenv("ADMIN_PASS") && user.Email == os.Getenv("ADMIN_EMAIL") {
		user.UserType = "ADMIN"
	} else {
		user.UserType = "USER"
	}

	if user.UserType == "ADMIN" {
		adminFilter := bson.M{"userType": "ADMIN"}
		if _, err := userCollection.FindOne(ctx, adminFilter).DecodeBytes(); err == nil {
			return c.Status(400).JSON(fiber.Map{
				"status":  "Fail",
				"message": "admin already exist",
				"data":    c.JSON(err),
			})
		}
	}

	emailFilter := bson.M{"email": user.Email}
	if _, err := userCollection.FindOne(ctx, emailFilter).DecodeBytes(); err == nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "user already exist",
			"data":    c.JSON(err),
		})
	}

	password, _ := utils.HashPassword(user.Password)

	user.Password = password

	if _, err := userCollection.InsertOne(ctx, user); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "fail",
			"message": "fail to insert into user",
			"data":    c.JSON(err),
		})
	}

	signedToken, err := utils.CreateToken(user.ID, user.Email, user.UserType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "fail",
			"message": "fail to create token",
			"data":    c.JSON(err),
		})
	}

	cookies := &fiber.Cookie{
		Name:     "jwt",
		Value:    signedToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}

	c.Cookie(cookies)

	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "User singed up successfull",
		"data":    user,
	})
}
