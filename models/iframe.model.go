package models

type Video struct {
	Download string `bson:"download" json:"download"`
	Iframe   string `bson:"iframe" json:"iframe"`
}
