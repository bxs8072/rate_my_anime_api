package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rating struct {
	Id        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID  `bson:"user" json:"user"`
	Anime     Anime               `bson:"anime" json:"anime"`
	Rating    float64             `bson:"rating" json:"rating"`
	CreatedAt primitive.Timestamp `bson:"createdAt" json:"createdAt"`
}
