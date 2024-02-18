package internal

var DRAFT string = "draft"
var PUBLISHED string = "published"

type Article struct {
	ID          int
	Title       string
	Markdown    string
	Description string
	Status      string
	Tags        []Tag
}

type Tag struct {
	ID   string
	Name string
}
