package handlebars

import (
	"gocms/models"

	"github.com/flowchartsman/handlebars/v3"
)

func Render(source string, data any) (string, error) {
	return handlebars.Render(source, data)
}

func NewView() models.View {
	return models.View{}
}
