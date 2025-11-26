package app

import (
	"fmt"
	"gocms/app/models"
	"gocms/templates"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type key string

const RequestIDKey key = "requestID"

type Request struct {
	Id          string
	View        models.View
	RouteType   string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
	Path        string
	UrlSegments []string
	AnyValues   []string
	NamedValues map[string]string
	Status      string
}

var Mutex sync.Mutex

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

// add views and render engine
var Requests = make(map[string]*Request)

var BaseUrl = ""

// voutputs["lol"] = "string"
//var wg *sync.WaitGroup

func AddView(w http.ResponseWriter, location string, args any) string {
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
			//head :=
			//fmt.Println(r)

			Requests[w.Header().Get("X-Request-Id")].View.Output = out

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
	return Requests["0"].View
}
func ClearOutput(id string) {
	Requests[string(id)].View.Output = ""
}

func GetOutput(w http.ResponseWriter) string {
	return Requests[w.Header().Get("X-Request-Id")].View.Output
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
