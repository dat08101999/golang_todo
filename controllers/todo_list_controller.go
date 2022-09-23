package controllers

import (
	"api_fiber/dbquery"
	"api_fiber/models"
	"api_fiber/res"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetListTodoJson(c *fiber.Ctx, userName string) error {
	todoList, err := dbquery.GetListTodo(userName)
	return res.Response(c, todoList, err, "Get list success")
}

func InsertTodo(c *fiber.Ctx, userName string) error {
	var listTodo struct {
		Todos []models.Todo `bson:"todos"`
	}
	c.BodyParser(&listTodo)
	var interfaceTodos []interface{} = make([]interface{}, len(listTodo.Todos))
	for i, d := range listTodo.Todos {
		d.Id = primitive.NewObjectID()
		d.Username = userName
		d.CreateAt = primitive.NewDateTimeFromTime(time.Now())
		interfaceTodos[i] = d
	}
	result, err := dbquery.InsertTodoListMany(interfaceTodos)
	return res.Response(c, nil, err, result)
}
