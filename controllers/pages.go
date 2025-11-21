package controller

import (
	"fmt"
	"gocms/app"
	"gocms/models"
	"net/http"
)

// controllers

func Index(welcome_message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		ctx := models.Post{
			models.Person{"Jean", "Valjean"},
			"Life is difficult",
			[]models.Comment{
				models.Comment{
					models.Person{"Marcel", "Beliveau"},
					"LOL!",
				},
			},
		}
		app.AddView("/home.hbs", ctx)
		fmt.Fprintf(w, "%v<a href='/test'>Test page</a>- <a href='/another/value1/and/value2'>Test page</a> ", string(welcome_message)+app.GetView().Output)
	})
}

func OtherHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Named Values %v\n </div>", app.Router.NamedValue("id"))
	}
	return http.HandlerFunc(fn)
}

func ServicesHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Any value 0: %v\n </div>", app.Router.AnyValue(0))
	}
	return http.HandlerFunc(fn)
}

func TestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		output := ""
		for _, value := range app.Router.AnyValues() {
			output += value
		}
		fmt.Fprintf(w, "<div>All any values %v\n </div>", output)
	}
	return http.HandlerFunc(fn)
}
