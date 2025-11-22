package view

import (
	"fmt"
	"gocms/models"
	"gocms/templates"
	"os"
	"path/filepath"
)

// add views and render engine

var v = models.View{}

func AddView(location string, args any) (string, error) {
	//no reason to choose engine for now with: app.GetConfig()...
	data, err := os.ReadFile("./views/" + location)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	err = nil

	if filepath.Ext(location) == ".hbs" {
		hb := templates.Handlebars{}
		v.Output, _ = hb.Render(string(data), args)
	} else if filepath.Ext(location) == ".mustache" {
		hb := templates.Mustache{}
		v.Output, _ = hb.Render(string(data), args)
	}
	//add more types of rendering here...
	return v.Output, nil
}

func GetView() models.View {
	return v
}
func ClearOutput() {
	v.Output = ""
}
