package mustache

import (
	"gocms/models"

	"github.com/cbroglie/mustache"
)

func Render(source string, data any) (string, error) {
	return mustache.Render(source, data)
}

func NewView() models.View {
	return models.View{}
}
