package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty" validate:"required"`
	Description string             `json:"description,omitempty" bson:"description,omitempty" validate:"required"`
	Category    string             `json:"category,omitempty" bson:"category,omitempty" validate:"required"`
	Brand       string             `json:"brand,omitempty" bson:"brand,omitempty" validate:"required"`
	Price       int                `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty" validate:"required"`
	Images      []image            `json:"images,omitempty" bson:"images,omitempty" validate:"required"`
}

type image struct {
	Image string `json:"image,omitempty" bson:"image,omitempty"`
}
