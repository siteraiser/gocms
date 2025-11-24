package autorouted

import (
	"fmt"
	"gocms/app"
	"net/http"
)

// auto-routed pages require this structure

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<div>Named Values %v\n </div>", app.NamedValues["id"])
}

func Page2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<div>Named Values %v\n </div>", app.NamedValues["id"])
}
