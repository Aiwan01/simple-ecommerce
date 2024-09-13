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

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

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

func Signin(c *fiber.Ctx) error {

	type SigninRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	ctx, err := context.WithTimeout(context.TODO(), 10*time.Second)
	defer err()

	var body SigninRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "error",
			"message": "some field is missing",
			"data":    err.Error(),
		})
	}

	var existingUser bson.Raw
	if err := userCollection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&existingUser); err != nil {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "User not exist",
			"data":    err.Error(),
		})
	}
	isValidPassword := utils.VerifyPassword(body.Password, existingUser.Lookup("password").StringValue())
	if !isValidPassword {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid credential",
			"data":    nil,
		})
	}
	// create token
	signinToken, errors := utils.CreateToken(existingUser.Lookup("_id").ObjectID(), existingUser.Lookup("email").StringValue(), existingUser.Lookup("userType").StringValue())

	if errors != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "Fail to create token",
			"data":    errors.Error(),
		})
	}

	cookies := &fiber.Cookie{
		Name:     "jwt",
		Value:    signinToken,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(cookies)

	return c.Status(500).JSON(fiber.Map{
		"status":  "success",
		"message": "User logged in successfully",
		"data":    existingUser,
	})
}

func Logout(c *fiber.Ctx) error {

	cookies := &fiber.Cookie{
		Name:    "jwt",
		Value:   "",
		Expires: time.Now().Add(-time.Hour),
	}

	c.Cookie(cookies)
	return c.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "successfully logged out",
		"data":    nil,
	})
}

func ProfileUser(c *fiber.Ctx) error {
	idLocal := c.Locals("id").(string)
	userId, err := primitive.ObjectIDFromHex(idLocal)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "fail to get user id",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	var user models.User
	err = userCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
			"data":    err.Error(),
		})
	}

	return c.Status(fiber.StatusFound).JSON(fiber.Map{
		"status":  "success",
		"message": "User fetched successfullu",
		"data":    user,
	})

}
