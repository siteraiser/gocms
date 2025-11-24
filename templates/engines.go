package templates

import (
	template "gocms/app/templates"
	"gocms/templates/engines/handlebars"
	"gocms/templates/engines/mustache"
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
	results, err := mustache.Render(s, a)
	return results, err
}

type Handlebars struct{}

func (h Handlebars) Render(s string, a any) (string, error) {
	results, err := handlebars.Render(s, a)
	return results, err
}
