package app

import (
	"gocms/app/router"
	"net/http"
)

type Routing struct {
	Handler http.Handler
}

var Router Routing

func NewApp(ah http.Handler) {
	Router = Routing{
		ah,
	}

}

func (h *Routing) AnyValue(index int) string {
	return router.AnyValue(index)
}
func (h *Routing) AnyValues() []string {
	return router.AnyValues()
}
func (h *Routing) NamedValue(name string) string {
	return router.NamedValue(name)
}

type Page struct {
	Attributes struct {
		Header http.Header
	}
	Meta    string
	Content string
	Assets  map[string]string
}
