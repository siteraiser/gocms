package templates

import (
	models "gocms/app/views/templates/defs"
)

var V models.View

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
