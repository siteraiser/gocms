package view

import (
	"github.com/flowchartsman/handlebars/v3"
)

func Render(source string, ctx any) string {
	return handlebars.MustRender(source, ctx)
}
