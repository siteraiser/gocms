package main

import (
	system "gocms/app/controllers"
	router "gocms/app/router"
	controller "gocms/controllers"
	"net/http"
)

func addRoutes() {
	for _, route := range []struct {
		pattern    string
		controller http.Handler
	}{
		{"/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets")))},
		{"/", controller.Index("Welcome!")},
		{"/testpage/{$}", controller.ServicesHandler("test")},
		{"GET /another/{$}/and/{$}", controller.TestHandler()},
		{"/another/{id}/link", controller.OtherHandler("test2")},
		{"/random", controller.RandoHandler()},
	} {
		router.Add(route.pattern, route.controller)
	}
}
