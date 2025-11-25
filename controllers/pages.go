package controller

import (
	"fmt"
	"gocms/app"
	"gocms/models"
	"math/rand"
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

		home := app.AddView("home.hbs", ctx)
		ctx2 := models.Page{
			Lang:  app.GetConfig().Settings.Preferences.Language,
			Title: string(welcome_message),
			Body:  home + "<a href='/test'>Test page</a><br><a href='/another/value1/and/value2'>Any Vars Test page</a><br><a href='/another/value1/link'>Named Vars Test page</a><br><a href='/autorouted'>Auto-routed page</a>",
		}
		fmt.Fprintf(w, "%v", app.AddView("document.hbs", ctx2))
	})
}

func RandoHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set headers to prevent caching
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		randomInt := rand.Intn(1000)
		fmt.Println("Random Integer:", randomInt)

		fmt.Fprintf(w, "<div>Randomness %v\n </div>", randomInt)
	}
	return http.HandlerFunc(fn)
}

func OtherHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Named Values %v\n </div>", app.NamedValues["id"])
	}
	return http.HandlerFunc(fn)
}

func TestPageHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>Any value 0: %v\n </div><br><a href='"+app.BaseUrl+"autorouted'>Go to autorouted page</a>", app.AnyValues[0])
	}
	return http.HandlerFunc(fn)
}

func TestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		output := ""
		for _, value := range app.AnyValues {
			output += value
		}
		fmt.Fprintf(w, "<div>All any values %v\n </div>", output)
	}
	return http.HandlerFunc(fn)
}
