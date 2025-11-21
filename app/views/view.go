package view

import (
	"fmt"

	"gocms/templates/handlebars"
	"gocms/templates/mustache"
	"os"
	"path/filepath"
)

// add views and render engine
type View struct {
	Args   any
	Output string
}

var v = View{}

func AddView(location string, args any) (string, error) {
	//no reason to choose engine for now with: app.GetConfig()...
	if filepath.Ext(location) == ".hbs" {
		data, err := os.ReadFile("./views/" + location)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		err = nil
		v.Output, err = handlebars.Render(string(data), args)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
	} else if filepath.Ext(location) == ".mustache" {
		data, err := os.ReadFile("./views/" + location)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
		err = nil
		v.Output, err = mustache.Render(string(data), args)
		if err != nil {
			fmt.Println("Error:", err)
			return "", err
		}
	}
	//add more types of rendering here...
	return v.Output, nil
}

func GetView() View {
	return v
}
func ClearOutput() {
	v.Output = ""
}

type Create interface {
	createView() View
}

type Mustache struct{}

func (m Mustache) createView() mustache.View {
	V := mustache.NewView()
	//V.Output
	return V
}

type Handlebars struct{}

func (h Handlebars) createView() handlebars.View {
	V := handlebars.NewView()
	return V
}
