package handlebars

import (
	"github.com/flowchartsman/handlebars/v3"
)

type View struct {
	Args   any
	Output string
}

func Render(source string, data any) (string, error) {

	return handlebars.Render(source, data)
}

func NewView() View {

	return View{}
}
