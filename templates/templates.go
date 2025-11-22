package templates

import (
	"fmt"
	"gocms/models"
	"gocms/templates/engines/handlebars"
	"gocms/templates/engines/mustache"
)

var V models.View

func RenderView(h any, s string, a any, renderer func(s string, a any) (string, error)) (string, error) {
	out, err := renderer(s, a)
	if err != nil {
		fmt.Println(err.Error())
	}
	V.Output = out
	return V.Output, err
}

type Engine struct {
	Ext    string
	Render Renderer
}
type Renderer interface {
	Render(string, any) (string, error)
}

type Engines struct {
	List []Engine
}

// Filled in by user...
var Registry = []Engine{
	hbs,
	mst,
}

var hbs = Engine{
	Ext:    ".hbs",
	Render: Handlebars{},
}

var mst = Engine{
	Ext:    ".mustache",
	Render: Mustache{},
}

type Mustache struct{}

func (h Mustache) Render(s string, a any) (string, error) {
	return RenderView(h, s, a, func(s string, a any) (string, error) {
		V = models.View{}
		results, err := mustache.Render(s, a)
		return results, err
	})
}

type Handlebars struct{}

func (h Handlebars) Render(s string, a any) (string, error) {
	return RenderView(h, s, a, func(s string, a any) (string, error) {
		V = models.View{}
		results, err := handlebars.Render(s, a)
		return results, err
	})
}
