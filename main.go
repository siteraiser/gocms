package main

import (
	"flag"
	"fmt"
	"gocms/app"
	"gocms/app/router"
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

	app.NewApp(&router.Handler{})

	addRoutes()

	fmt.Println("Starting server on :" + port)
	if err := http.Serve(listener, app.Router.Handler); err != nil {
		panic(err)
	}
}
