package models

import (
	db "gocms/app/database"
	"strings"
)

func LoadByLink(link string) string {
	rows, _ := db.Db.Query("SELECT content FROM pages WHERE link = ?", strings.TrimLeft(link, "/"))
	var content string
	type page struct {
		Content string
	}
	p := page{}
	for rows.Next() {
		rows.Scan(&content)
		p.Content = content
	}
	return p.Content
}
