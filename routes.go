package main

import (
	system "gocms/app/controllers"
	router "gocms/app/router"
	controller "gocms/controllers"
	. "gocms/controllers/autorouted"
	"net/http"
)

func addRoutes() {
	router.Add("/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets"))))
	router.Add("/", controller.Index("Welcome to GoCMS!"))
	router.Add("/blog", controller.BlogHandler(""))
	router.Add("/testpage/{$}", controller.TestPageHandler("test"))
	router.Add("GET /another/{$}/and/{$}", controller.TestHandler())
	router.Add("/another/{id}/link", controller.OtherHandler("test2"))
	router.Add("/random", controller.RandoHandler())
	router.Add("/form-submissions", controller.FormHandler())

	//auto-mvc functions
	router.AddFunc(Index)
	router.AddFunc(Page2)
	router.AddFunc(Page3)
}
