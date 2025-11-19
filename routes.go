package main

import (
	system "gocms/app/controllers"
	router "gocms/app/router"
	controller "gocms/controllers"
	"net/http"
)

func addRoutes() {
	router.Add("/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets"))))
	router.Add("/", controller.Index("Welcome!"))
	router.Add("/testpage/{$}", controller.ServicesHandler("test"))
	router.Add("/another/{$}/and/{$}", controller.TestHandler())
	router.Add("/another/{id}/link", controller.OtherHandler("test2"))
}
