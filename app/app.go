package app

import (
	"fmt"
	"gocms/app/router"
	view "gocms/app/viewers"
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

func (h *Routing) AnyValue(index int) string {
	return router.AnyValue(index)
}
func (h *Routing) AnyValues() []string {
	return router.AnyValues()
}
func (h *Routing) NamedValue(name string) string {
	return router.NamedValue(name)
}

// add views and render engine
type View struct {
	Content  string
	Location string
	Renderer any
	Args     any
}

type Renderer interface {
	Show() string
}

var v = View{}

func AddView(location string, args any) error {

	if filepath.Ext(location) == ".hbs" {
		data, err := os.ReadFile(location)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		v.Content = view.Show(string(data), args)
	}
	//add more types of rendering here...
	return nil
}
func GetContent() string {
	return v.Content
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
