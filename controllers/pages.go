package controller

import (
	"fmt"
	"gocms/router"
	"net/http"
)

// controllers
func OtherHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Named Values %v\n </div>", router.NamedValue("id"))
	}
	return http.HandlerFunc(fn)
}

func ServicesHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Any value 0: %v\n </div>", router.AnyValue(0))
	}
	return http.HandlerFunc(fn)
}
