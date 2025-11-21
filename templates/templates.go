package templates

import (
	"gocms/models"
	"gocms/templates/engines/handlebars"
	"gocms/templates/engines/mustache"
)

var V models.View

type Mustache struct{}

func (m Mustache) CreateView() {
	V = mustache.NewView()
}
func (h Mustache) Render(s string, a any) {
	V.Output, _ = mustache.Render(s, a)
}

type Handlebars struct{}

func (h Handlebars) CreateView() {
	V = handlebars.NewView()
}

func (h Handlebars) Render(s string, a any) {
	V.Output, _ = handlebars.Render(s, a)
}
