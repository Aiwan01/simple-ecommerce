package routes

import (
	"go-ecom/controllers"

	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {

	api := app.Group("/api/v1")
	api.Get("/", controllers.Welcome)

	//auth router
	userApi := app.Group("/user/auth")
	userApi.Post("/signup")

}
