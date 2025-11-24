package autorouted

import (
	"fmt"
	"gocms/app"
	"gocms/models"
	"net/http"
)

// auto-routed pages require this structure

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx := models.Page{
		Title: "Hello, World!",
	}

	index, _ := app.AddView("index.mustache", ctx)
	app.ClearOutput()

	app.AddView("layouts/main.mustache", map[string]string{"embed": index})

	fmt.Fprintf(w, "%v", app.GetOutput())
}

func Page2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<div>autoloaded.Page2  </div>")
}
