package handlebars

import (
	"github.com/flowchartsman/handlebars/v3"
)

func Render(source string, data any) (string, error) {
	return handlebars.Render(source, data)
}
