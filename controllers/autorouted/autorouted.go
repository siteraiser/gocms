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

	embed := app.Render(
		"index.mustache",
		map[string]any{
			"Title": title,
			"Linklist": []map[string]interface{}{
				{"Text": "Home", "Link": "/"},
				{"Text": "Random", "Link": "/random"},
				{"Text": "Test page 2", "Link": "/autorouted/page2"},
				{"Text": "Test POST form", "Link": "/form-submissions"},
			},
		},
	)

	page := app.Render(
		"layouts/main.mustache",
		map[string]string{
			"Lang":  app.Config.Settings.Language,
			"Title": title,
			"embed": embed,
		},
	)

	//	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "%v", page)
}

func Page2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	cms := app.Cms(r)
	cms.Views.Add(
		"layouts/main.mustache",
		map[string]string{
			"Lang":  app.Config.Settings.Language,
			"Title": "Page 2",
			"embed": "<div>autorouted.Page2  </div><a href=\"/autorouted/page3\">Page 3</a>",
		})
	time.Sleep(time.Second * 4)
}

func Page3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	App := app.Cms(r)
	time.Sleep(time.Second * 2)
	App.Views.Add(
		"layouts/main.mustache",
		map[string]string{
			"Lang":  app.Config.Settings.Language,
			"Title": "Page 3",
			"embed": "Some Other Contents<div>autorouted.Page3  </div",
		})
	time.Sleep(time.Second * 2)
}
