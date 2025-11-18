package main

import (
	"flag"
	"fmt"

	controller "gocms/controllers"
	system "gocms/controllers/system"
	"gocms/router"
	"net"
	"net/http"
	"strconv"
)

func main() {

	portFlag := flag.Int("port", 8080, "string")
	flag.Parse()
	port := strconv.Itoa(*portFlag)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}
	//Serve assets
	//http.Handle("/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets"))))
	//Add routes
	router.Add("/assets/", http.StripPrefix("/assets/", system.Fs(http.Dir("./assets"))))
	router.Add("/testpage/{$}", controller.ServicesHandler("test"))
	//router.Add("/another/{$}/and/{$}", controller.TestHandler())
	router.Add("/another/{id}/link", controller.OtherHandler("test2"))

	handler := &router.Handler{}
	fmt.Println("Starting server on :8080")
	if err := http.Serve(listener, handler); err != nil {
		panic(err)
	}
}
