package controller

import (
	"fmt"
	"gocms/app"
	"gocms/app/sys"
	"gocms/models"
	samplemodels "gocms/models/cms"
	"gocms/modules/menus"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
)

// controllers
func Index(welcome_message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cms := app.Cms(r)

		page := cms.Utils.NewPage()
		page.Title = welcome_message
		page.Body = strconv.Itoa(sys.Stats.TotalHits) +
			" page" +
			cms.Utils.Grammar.LowerIfPluralS(sys.Stats.TotalHits) +
			" served this session."

		cms.Views.Add("document.hbs", page)
	})
}

// controllers
func BlogHandler(welcome_message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cms := app.Cms(r)

		post := models.Post{
			Author: models.Person{FirstName: "Jean", LastName: "Valjean"},
			Body:   "Life is difficult",
			Comments: []models.Comment{
				{
					Author: models.Person{FirstName: "Marcel", LastName: "Beliveau"},
					Body:   "LOL!",
				},
			},
		}
		home := cms.Views.Render("posttemplate.hbs", post)

		page := cms.Utils.NewPage()
		page.Title = "Welcome to the Handlebars example blog template"
		page.Body = home + "<a href='/'>Home</a>"

		cms.Views.Add("document.hbs", page)
	})
}

func RandoHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Set headers to prevent caching
		cms := app.Cms(r).Header.Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		cms.Header.Set("Pragma", "no-cache")
		cms.Header.Set("Expires", "0")
		//cms.Header.Set("Content-Type", "text/html; charset=utf-8")

		worker := func(s *string, wg *sync.WaitGroup) {
			defer wg.Done()
			*s += strconv.Itoa(rand.Intn(1000000000000000))
		}

		random := ""
		var wg sync.WaitGroup

		// Launch several workers
		for i := 1; i <= 100; i++ {
			wg.Add(1) // Increment the WaitGroup counter
			go worker(&random, &wg)
		}
		wg.Wait()
		random = "<p style=\"word-break: break-all;\">" + random + "</p>"
		//fmt.Println("Random Integer:", random)

		cms.Views.SetOut(cms.Views.Render("document.hbs", models.Page{Body: random}))

	}
	return http.HandlerFunc(fn)
}

func TestPageHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprintf(w, "<div>Any value 0: %v\n </div><br><a href='"+app.BaseUrl+"autorouted'>Go to autorouted page</a>", app.Cms(r).Any.Values()) //app.Req(r).AnyValues[0]
	}
	return http.HandlerFunc(fn)
}

func TestHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//App := app.Cms(r)
		output := ""
		for _, value := range app.Cms(r).Any.Vals {
			output += value
		}
		fmt.Fprintf(w, "<div>All any values %v\n </div>", output)
	}
	return http.HandlerFunc(fn)
}

func OtherHandler(arguments string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//app.NamedValues(r)["id"]
		//app.Req(r).NamedValues["id"]
		//or...for now...
		cms := app.Cms(r)

		header := app.Render("partials/header.hbs", nil)

		//time.Sleep(time.Second * 3)
		link := cms.URL.Path() + "?params1[]=value1&params1[]=value2&param2=value1"
		output := "<div>Name:  - Value: </div>"
		for name, value := range cms.Named.Vals {
			output += "<div>" + string(name) + " - " + string(value) + "</div>"
		}

		fmt.Fprintf(
			w,
			header+cms.Utils.Html.Ahref(link, "With params")+"<div>Values:<b>%v</b> </div>"+output,
			cms.URL.QueryParams(),
		)
	}
	return http.HandlerFunc(fn)
}

func FormHandler() http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//app.NamedValues(r)["id"]
		//app.Req(r).NamedValues["id"]
		//or...for now...
		cms := app.Cms(r)
		input := cms.Form.Fields().Get("posted_value")
		samplemodels.UpdatePage(input)
		content := cms.Views.Render("forms/testform.hbs", nil)

		fmt.Fprintf(
			w,
			menus.Nav()+"<h2>Update the <a href=\"/cms-page\">CMS Page</a></h2>"+content+"<div>Values:<b>%v</b> </div>",
			"Page Updated.",
		)
	}
	return http.HandlerFunc(fn)
}
