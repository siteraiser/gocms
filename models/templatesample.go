package models

type Person struct {
	FirstName string
	LastName  string
}

type Comment struct {
	Author Person
	Body   string `handlebars:"content"`
}

type Post struct {
	Author   Person
	Body     string
	Comments []Comment
}

type Page struct {
	Lang    string
	Welcome string
	Body    string
}
