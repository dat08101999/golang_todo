package dbquery

import (
	"api_fiber/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertUser(user models.User) error {
	var isRegisterd error
	query(func(client *mongo.Client, ctx context.Context) {
		colection := getColecttion(client, "user")
		var result bson.M
		errFindOne := colection.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&result)
		if errFindOne == nil {
			isRegisterd = errors.New("user dose exist")
		} else {
			_, err2 := colection.InsertOne(ctx, user)
			if err2 != nil {
				isRegisterd = err2
			}
		}
	})
	return isRegisterd
}

func FindUser(user models.User) error {
	var isFindUser error
	query(func(client *mongo.Client, ctx context.Context) {
		colection := getColecttion(client, "user")
		var result bson.M
		errFindOne := colection.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&result)
		if errFindOne != nil {
			isFindUser = errors.New("User name does not exist")
		} else {
			hashedPassword := fmt.Sprintf("%v", result["password"])
			err2 := models.CheckPasswordHash(hashedPassword, user.Password)
			if err2 != nil {
				isFindUser = errors.New("Password is incorrect")
			}
		}
	})
	return isFindUser
}
