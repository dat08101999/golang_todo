package dbquery

import (
	"api_fiber/models"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func insertOne() {
	query(func(client *mongo.Client, ctx context.Context) {

	})
}

func deleteMany() {
	query(func(client *mongo.Client, ctx context.Context) {

	})
}

func deleteOne() {
	query(func(client *mongo.Client, ctx context.Context) {

	})
}

func InsertTodoListMany(todoList []interface{}) (string, error) {
	var err error
	query(func(client *mongo.Client, ctx context.Context) {
		collection := getColecttion(client, "todo")
		_, errInsert := collection.InsertMany(ctx, todoList, options.InsertMany())
		if errInsert != nil {
			err = errInsert
		}
	})
	return "success", err
}

func GetListTodo(userName string) (todolist []models.Todo, err error) {
	var errorG error
	var todoList []models.Todo
	query(func(client *mongo.Client, ctx context.Context) {
		collection := getColecttion(client, "todo")
		cur, err := collection.Find(ctx, bson.M{
			"username": userName,
		})
		if err != nil {
			errorG = err
		}
		for cur.Next(ctx) {
			var t models.Todo
			err := cur.Decode(&t)
			if err != nil {
				println("error")
			}
			todoList = append(todoList, t)
		}

	})

	return todoList, errorG
}
