package routes

import (
	"go-ecom/controllers"
	"go-ecom/middleware"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {

	api := app.Group("/api/v1")
	api.Get("/", controllers.Welcome)

	//auth router
	userApi := app.Group("/user/auth")
	userApi.Post("/signup", middleware.RegisterCredentialInput, controllers.Signup)
	userApi.Post("/signin", middleware.RegisterCredentialInput, controllers.Signin)
	userApi.Post("/signout", middleware.RequireAuthValidate, controllers.Logout)
	userApi.Get("/profile", middleware.RequireAuthValidate, controllers.ProfileUser)

	productApi := app.Group("/products", middleware.RequireAuthValidate)
	productApi.Post("/create", controllers.CreateNewProducts)
	productApi.Get("/", controllers.GetAllProducts)
	productApi.Get("/:id", controllers.GetProduct)
	productApi.Put("/:id", controllers.UpdateProduct)
	productApi.Delete("/:id", controllers.DeleteProduct)

}
