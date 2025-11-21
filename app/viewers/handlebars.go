package view

import (
	"github.com/flowchartsman/handlebars/v3"
)

func Show(source string, ctx any) string {
	return handlebars.MustRender(source, ctx)
}
