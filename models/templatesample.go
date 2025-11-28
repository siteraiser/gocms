package models

type Page struct {
	Lang  string
	Title string
	Body  string
}

type Link struct {
	Href string
	Text string
}
type TestIndexPage struct {
	Lang     string
	Title    string
	Linklist []Link
}

/*
	type Page struct {
		Attributes struct {
			Header http.Header
		}
		Meta    string
		Content string
		Assets  map[string]string
	}
*/
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
