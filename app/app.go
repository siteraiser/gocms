package app

import (
	"gocms/app/router"
	view "gocms/app/views"
	"net/http"
)

// add routing from routing package to app
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

// add views
func AddView(location string, args any) error {
	return view.AddView(location, args)
}

func GetView() view.View {
	return view.GetView()
}

func GetOutput() string {
	return view.GetView().Output
}

func ClearOutput() {
	view.ClearOutput()
}

/*
type Page struct {
	Attributes struct {
		Header http.Header
	}
	Meta    string
	Content string
	Assets  map[string]string
}
*/
