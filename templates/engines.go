package templates

import (
	template "gocms/app/templates"

	"github.com/cbroglie/mustache"
	"github.com/flowchartsman/handlebars/v3"
)

// Filled in by user...
var Registry = []template.Engine{
	{
		Ext:    ".hbs",
		Engine: Handlebars{},
	},
	{
		Ext:    ".mustache",
		Engine: Mustache{},
	},
}

type Mustache struct{}

func (h Mustache) Render(s string, a any) (string, error) {
	return mustache.Render(s, a)

}

type Handlebars struct{}

func (h Handlebars) Render(s string, a any) (string, error) {
	return handlebars.Render(s, a)
}
