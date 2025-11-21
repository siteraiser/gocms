package view

import (
	"fmt"
	"os"
	"path/filepath"
)

// add views and render engine
type View struct {
	Args   any
	Output string
}

var v = View{}

func AddView(location string, args any) error {

	if filepath.Ext(location) == ".hbs" {
		data, err := os.ReadFile("./views/" + location)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		v.Output = Render(string(data), args)

	}
	//add more types of rendering here...
	return nil
}

func GetView() View {
	return v
}
func ClearOutput() {
	v.Output = ""
}
