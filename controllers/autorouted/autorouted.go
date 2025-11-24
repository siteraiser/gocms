package autorouted

import (
	"fmt"
	"gocms/app"
	"net/http"
)

// auto-routed pages require this structure

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	ctx := map[string]any{
		"Lang":  app.GetConfig().Settings.Preferences.Language,
		"Title": "Another",
		"Linklist": []map[string]interface{}{
			{"Text": "Home", "Link": "/"},
			{"Text": "Random", "Link": "/random"},
			{"Text": "Test page 2", "Link": "/autorouted/page2"},
		},
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
