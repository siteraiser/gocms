package mustache

import (
	"github.com/cbroglie/mustache"
)

func Render(source string, data any) (string, error) {
	return mustache.Render(source, data)
}
