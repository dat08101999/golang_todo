package controllers

import (
	"api_fiber/dbquery"
	"api_fiber/res"

	"github.com/gofiber/fiber/v2"
)

func GetListTodoJson(c *fiber.Ctx) error {
	todoList, err := dbquery.GetListTodo()
	return res.Response(c, todoList, err, "Get list success")
}
