package app

import (
	view "gocms/app/views"
	models "gocms/app/views/templates/defs"
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

var BaseUrl = ""             // Is set during the installation in config. Should always contain a trailing slash.
var Path = ""                // Contains full current path.
var UrlSegments = []string{} // Contains the url segements for the current page. Access the first part of the path using app.url_segments[0].
var Request *http.Request
var AnyValues = []string{}
var NamedValues = map[string]string{}
var RouteType = ""

/*
// Parameter Functions
func AnyValues() []string {
	return AnyParams
}
func AnyValue(index int) string {

	if len(AnyParams)-1 >= index {
		return AnyParams[index]
	}
	return ""
}
*/

// add views
func AddView(location string, args any) (string, error) {
	return view.AddView(location, args)
}

func GetView() models.View {
	return view.GetView()
}

func GetOutput() string {
	res := view.GetView()
	return res.Output
}

func ClearOutput() {
	view.ClearOutput()
}
