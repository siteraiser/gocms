package main

import (
	controller "gocms/controllers"
	system "gocms/core/controllers"
	"gocms/core/router"
	"net/http"
)

func addRoutes() {
	router.Add("/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets"))))
	router.Add("/", controller.Index("Welcome!"))
	router.Add("/testpage/{$}", controller.ServicesHandler("test"))
	router.Add("/another/{$}/and/{$}", controller.TestHandler())
	router.Add("/another/{id}/link", controller.OtherHandler("test2"))
}
