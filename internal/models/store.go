package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Store struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	ProfilePic  string             `json:"profile_pic,omitempty" bson:"profile_pic,omitempty" validate:"required"`
	Products    []product          `json:"products,omitempty" bson:"products,omitempty"`
	CreatedAt   string             `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type product struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedAt string             `json:"created_at,omitempty" bson:"created_at,omitempty"`
}
