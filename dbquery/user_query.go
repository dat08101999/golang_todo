package dbquery

import (
	"api_fiber/models"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			user.Id = primitive.NewObjectID()
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

func GetUser(userName string) *models.User {
	var user *models.User
	query(
		func(client *mongo.Client, ctx context.Context) {
			collecttion := getColecttion(client, "user")
			collecttion.FindOne(ctx, bson.M{"username": userName}).Decode(&user)
		})
	return user
}
func UpdateUser(user models.User) error {
	var isError error
	query(func(client *mongo.Client, ctx context.Context) {
		if FindUser(user) == nil {
			isError = errors.New("User does not exist")
		} else {
			collecttion := getColecttion(client, "user")

			result, errorUpdate := collecttion.UpdateOne(ctx, bson.D{{Key: "username", Value: user.UserName}}, bson.D{{Key: "$set",
				Value: bson.D{
					{"refresh_token", user.Refresh_Token},
				},
			}}, &options.UpdateOptions{})
			if errorUpdate != nil {
				isError = errorUpdate
			} else {
				println(result)
			}
		}
	})
	return isError
}
