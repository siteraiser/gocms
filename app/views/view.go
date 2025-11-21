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

		templates.Handlebars.CreateView(templates.Handlebars{})
		templates.Handlebars.Render(templates.Handlebars{}, string(data), args)
		v.Output = templates.V.Output

	} else if filepath.Ext(location) == ".mustache" {

		templates.Mustache.CreateView(templates.Mustache{})
		templates.Mustache.Render(templates.Mustache{}, string(data), args)
		v.Output = templates.V.Output

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

/**/
