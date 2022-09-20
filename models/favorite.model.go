package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Favorite struct {
	Id        primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID  `bson:"user" json:"user"`
	Anime     Anime               `bson:"anime" json:"anime"`
	Status    string              `bson:"status" json:"status"`
	CreatedAt primitive.Timestamp `bson:"createdAt" json:"createdAt"`
}
