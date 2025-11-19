package controller

import (
	"fmt"
	"gocms/router"
	"net/http"
)

// controllers

func Index(welcome_message string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "%v<a href='/another/value1/and/value2'>Test page</a>", string(welcome_message))
	}
	return http.HandlerFunc(fn)
}

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

func TestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		output := ""
		for _, value := range router.AnyValues() {
			output += value
		}
		fmt.Fprintf(w, "<div>All any values %v\n </div>", output)
	}
	return http.HandlerFunc(fn)
}
