package mustache

import (
	"github.com/cbroglie/mustache"
)

type View struct {
	Args   any
	Output string
}

func Render(source string, data any) (string, error) {

	return mustache.Render(source, data)
}

func NewView() View {
	return View{}
}
