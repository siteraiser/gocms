package menus

import (
	"encoding/json"
	"fmt"
)

type Menu struct {
	Items []Item
}
type Item struct {
	Name string `json:"name"`
	Url  string `json:"url"`

	Children []Item `json:"children"`
}

func Nav() string {
	menujsonstr := `[{
	"url":"","name":"Main Menu",
	"children":[
		{"url":"/","name":"Home"},
		{"url":"/autorouted","name":"Auto-routed Index","children":[{"url":"/autorouted/page2","name":"Autorouted page2"}]},
		{"url":"/test","name":"Test Page"}
		]
	},{
	"url":"","name":"Other Section",
	"children":[
		{"url":"/another/value1/and/value2","name":"Any Vars Test page"},
		{"url":"/another/value1/link","name":"Named Vars Test Page"}
		]
	}]`
	// Unmarshal the JSON data
	var menu Menu
	err := json.Unmarshal([]byte(menujsonstr), &menu.Items)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	var output = ""
	MakeMenu(menu.Items, &output)
	return "<ul>" + output + "</ul>"
}
func MakeMenu(items []Item, output *string) ([]Item, *string) {
	for _, item := range items {
		href := ""
		if item.Url != "" {
			href = "<a class=\"nav-link\" href=\"" + item.Url + "\">" + item.Name + "</a>"
		} else {
			href = "<span>" + item.Name + "</span>"
		}
		*output = *output + "<li class=\"nav-item\">" + href
		if len(item.Children) > 0 {
			*output += "<ul>"
			MakeMenu(item.Children, output)
			*output += "</ul>"
		}
		*output += "</li>"
	}
	return []Item{}, output
}
func FooterNav() string {
	menujsonstr := `[{
	"url":"","name":"Main Menu",
	"children":[
		{"url":"/","name":"Home"},
		{"url":"/autorouted","name":"Auto-routed Index","children":[{"url":"/autorouted/page2","name":"Autorouted page2"}]},
		{"url":"/test","name":"Test Page"}
		]
	},{
	"url":"","name":"Other Section",
	"children":[
		{"url":"/another/value1/and/value2","name":"Any Vars Test page"},
		{"url":"/another/value1/link","name":"Named Vars Test Page"}
		]
	}]`
	// Unmarshal the JSON data
	var menu Menu
	err := json.Unmarshal([]byte(menujsonstr), &menu.Items)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
	}
	var output = ""
	MakeMenu(menu.Items, &output)
	return "<ul>" + output + "</ul>"
}
