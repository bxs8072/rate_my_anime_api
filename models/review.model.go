package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Review struct {
	Id        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID   `bson:"user" json:"user"`
	Anime     Anime                `bson:"anime" json:"anime"`
	Message   string               `bson:"message" json:"message"`
	CreatedAt primitive.Timestamp  `bson:"createdAt" json:"createdAt"`
	Likes     []primitive.ObjectID `bson:"likes" json:"likes"`
	Dislikes  []primitive.ObjectID `bson:"dislikes" json:"dislikes"`
}
