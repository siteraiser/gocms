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

	for i, e := range templates.Registry.List {
		if e.Ext == filepath.Ext(location) {
			v.Output, _ = templates.Registry.List[i].Render.Render(string(data), args)
		}
	}
	/*
		if filepath.Ext(location) == ".hbs" {

		} else if filepath.Ext(location) == ".mustache" {
			hb := templates.Mustache{}
			v.Output, _ = hb.Render(string(data), args)
		}
	*/

	//add more types of rendering here...
	return v.Output, nil
}

func GetView() models.View {
	return v
}
func ClearOutput() {
	v.Output = ""
}
