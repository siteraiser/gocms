package main

import (
	"flag"
	"fmt"
	"gocms/app"
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

	App := app.NewApp(&app.Handler{})

	//Add routes	//Serve assets
	addRoutes()

	fmt.Println("Starting server on :" + port)
	if err := http.Serve(listener, App.Router.Handler); err != nil {
		panic(err)
	}
}
