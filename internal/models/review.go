package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Review struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductID primitive.ObjectID `json:"product_id,omitempty" bson:"product_id,omitempty" validate:"required"`
	Reviews   []UserReview       `json:"reviews,omitempty" bson:"reviews,omitempty" validate:"required"`
}

type UserReview struct {
	User      UserType `json:"user,omitempty" bson:"user,omitempty"`
	Comment   string   `json:"comment,omitempty" bson:"comment,omitempty" validate:"required"`
	Rating    int      `json:"rating,omitempty" bson:"rating,omitempty" validate:"required"`
	CreatedAt string   `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type UserType struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string             `json:"name,omitempty" bson:"name,omitempty"`
}
