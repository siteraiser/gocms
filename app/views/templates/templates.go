package templates

import (
	"fmt"
	models "gocms/app/views/templates/defs"
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
