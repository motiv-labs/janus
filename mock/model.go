package mock

type Tag string

type Recipe struct {
	Name string `bson:"name" json:"name"`
	Tags []Tag  `bson:"tags" json:"tags"`
}
