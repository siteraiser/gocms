package autorouted

import (
	"fmt"
	"gocms/app"
	"net/http"
	"time"
)

// auto-routed pages require this structure

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	title := "Another Title"

	embed := app.AddView(
		"index.mustache",
		map[string]any{
			"Title": title,
			"Linklist": []map[string]interface{}{
				{"Text": "Home", "Link": "/"},
				{"Text": "Random", "Link": "/random"},
				{"Text": "Test page 2", "Link": "/autorouted/page2"},
			},
		},
	)

	page := app.AddView(
		"layouts/main.mustache",
		map[string]string{
			"Lang":  app.Config.Settings.Preferences.Language,
			"Title": title,
			"embed": embed,
		},
	)

	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "%v", page)
}

func Page2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<div>autoloaded.Page2  </div>")
}
