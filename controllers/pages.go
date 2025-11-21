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

		home, _ := app.AddView("home.hbs", ctx)
		app.ClearOutput()

		ctx2 := models.Page{
			Lang:    app.GetConfig().Settings.Preferences.Language,
			Welcome: string(welcome_message),
			Body:    home + "<a href='/test'>Test page</a><br><a href='/another/value1/and/value2'>Any Vars Test page</a><br><a href='/another/value1/link'>Named Vars Test page</a>",
		}

		app.AddView("document.hbs", ctx2)
		fmt.Fprintf(w, "%v", app.GetOutput())
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
