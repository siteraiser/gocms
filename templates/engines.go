package templates

import (
	template "gocms/app/views/templates"
	models "gocms/app/views/templates/defs"
	"gocms/templates/engines/handlebars"
	"gocms/templates/engines/mustache"
)

// Filled in by user...
var Registry = []template.Engine{
	{
		Ext:    ".hbs",
		Render: Handlebars{},
	},
	{
		Ext:    ".mustache",
		Render: Mustache{},
	},
}

type Mustache struct{}

func (h Mustache) Render(s string, a any) (string, error) {
	return template.RenderView(h, s, a, func(s string, a any) (string, error) {
		template.V = models.View{}
		results, err := mustache.Render(s, a)
		return results, err
	})
}

type Handlebars struct{}

func (h Handlebars) Render(s string, a any) (string, error) {
	return template.RenderView(h, s, a, func(s string, a any) (string, error) {
		template.V = models.View{}
		results, err := handlebars.Render(s, a)
		return results, err
	})
}
