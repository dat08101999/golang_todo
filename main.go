package main

import (
	"api_fiber/controllers"
	"api_fiber/middlewares"
	"api_fiber/routers"

	_ "github.com/arsmn/fiber-swagger/v2/example/docs"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Get(routers.GetListTodo, middlewares.CheckToken(controllers.GetListTodoJson))
	app.Post("/User/Register", controllers.RegisterUser)
	app.Post("/User/Login", controllers.LoginUser)
	app.Post("/User/RefreshToken", controllers.RefreshToken)
	app.Listen(":8080")
}
