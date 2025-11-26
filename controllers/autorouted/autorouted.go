package autorouted

import (
	"fmt"
	"gocms/app"
	"net/http"
	"time"
)

// auto-routed pages require this structure

func Index(w http.ResponseWriter, r *http.Request) {
	head := r.Context()
	fmt.Println(head.Value(app.RequestIDKey))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	title := "Another Title"
	ctx := map[string]any{
		"Title": title,
		"Linklist": []map[string]interface{}{
			{"Text": "Home", "Link": "/"},
			{"Text": "Random", "Link": "/random"},
			{"Text": "Test page 2", "Link": "/autorouted/page2"},
		},
	}

	app.AddView(w,
		"layouts/main.mustache",
		map[string]string{
			"Lang":  app.Config.Settings.Preferences.Language,
			"Title": title,
			"embed": app.AddView(w,
				"index.mustache",
				ctx,
			),
		},
	)
	time.Sleep(time.Second * 5)
	fmt.Fprintf(w, "%v", app.GetOutput(w))
}

func Page2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<div>autoloaded.Page2  </div>")
}
