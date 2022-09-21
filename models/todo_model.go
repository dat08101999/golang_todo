package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `json:"name"`
	CreateAt  primitive.DateTime `json:"createAt"`
	ExpiredAt primitive.DateTime `json:"expiredAt"`
}
