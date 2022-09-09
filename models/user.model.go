package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UID          string             `bson:"uid" json:"uid"`
	Email        string             `bson:"email" json:"email"`
	FirstName    string             `bson:"firstName" json:"firstName"`
	MiddleName   string             `bson:"middleName" json:"middleName"`
	LastName     string             `bson:"lastName" json:"lastName"`
	DisplayImage string             `bson:"displayImage" json:"displayImage"`
}
