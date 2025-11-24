package templates

import (
	models "gocms/app/templates/defs"
)

var V models.View

type Engine struct {
	Ext    string
	Engine Renderer
}
type Renderer interface {
	Render(string, any) (string, error)
}

type Engines struct {
	List []Engine
}
