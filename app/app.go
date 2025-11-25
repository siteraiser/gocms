package app

import (
	"fmt"
	"gocms/app/models"
	"gocms/templates"
	"net/http"
	"os"
	"path/filepath"
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

// add views and render engine

var v = models.View{}

func AddView(location string, args any) string {
	out := ""
	//no reason to choose engine for now with: app.GetConfig()...
	data, err := os.ReadFile("./views/" + location)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	err = nil
	//Find the rendering engine in the registry (outside of app folder) and render
	for _, e := range templates.Registry {
		if e.Ext == filepath.Ext(location) {
			out, err = e.Engine.Render(string(data), args)
			v.Output = out
			if err != nil {
				fmt.Println("Error:", err)
				return ""
			}
		}
	}

	//add more types of rendering here...
	return out
}

func GetView() models.View {
	return v
}
func ClearOutput() {
	v.Output = ""
}

func GetOutput() string {
	return v.Output
}

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


func RenderView(h any, s string, a any, renderer func(s string, a any) (string, error)) (string, error) {
	out, err := renderer(s, a)
	if err != nil {
		fmt.Println(err.Error())
	}
	V.Output = out
	return V.Output, err
}




*/
