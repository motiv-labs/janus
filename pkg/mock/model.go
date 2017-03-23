package mock

// Tag represents the recipe tags
type Tag string

// Recipe represents a hellofresh recipe
type Recipe struct {
	Name string `bson:"name" json:"name"`
	Tags []Tag  `bson:"tags" json:"tags"`
}
