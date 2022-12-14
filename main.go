package main

import (
	"api_fiber/controllers"
	"api_fiber/middlewares"
	"api_fiber/routers"
	"os"

	_ "github.com/arsmn/fiber-swagger/v2/example/docs"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	app.Post("Todo/Insert", middlewares.CheckToken(controllers.InsertTodo))
	app.Get(routers.GetListTodo, middlewares.CheckToken(controllers.GetListTodoJson))
	app.Post("/User/Register", controllers.RegisterUser)
	app.Post("/User/Login", controllers.LoginUser)
	app.Post("/User/RefreshToken", controllers.RefreshToken)
	app.Post("/Todo/Delete", middlewares.CheckToken(controllers.DeleteTodo))
	// app.Listen(":8080")
	app.Listen(":" + os.Getenv("PORT"))

}
