package models

type AnimeDetail struct {
	Type     string    `bson:"type" json:"type"`
	Plot     string    `bson:"plot" json:"plot"`
	Genres   []string  `bson:"genres" json:"genres"`
	Released string    `bson:"released" json:"released"`
	Status   string    `bson:"status" json:"status"`
	Other    string    `bson:"other" json:"other"`
	Episodes []Episode `bson:"episodes" json:"episodes"`
}

type Episode struct {
	Title string `bson:"title" json:"title"`
	Link  string `bson:"link" json:"link"`
}
