package main

import (
	"flag"
	"fmt"
	"gocms/app"
	"gocms/app/router"
	"log"
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
	_, err = listener.Accept()
	if err != nil {
		log.Printf("Accept error: %v", err)

	}

	app.NewApp(&router.Handler{})

	//Add routes	//Serve assets
	addRoutes()

	fmt.Println("Starting server on :" + port)
	if err := http.Serve(listener, app.Router.Handler); err != nil {
		panic(err)
	}
}
