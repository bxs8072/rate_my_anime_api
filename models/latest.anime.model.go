package models

type LatestAnime struct {
	Title string `bson:"title" json:"title"`
	Image string `bson:"image" json:"image"`
	Link  string `bson:"link" json:"link"`
}
