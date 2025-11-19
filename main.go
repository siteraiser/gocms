package main

import (
	"flag"
	"fmt"
	"gocms/core/router"
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

	//Add routes
	addRoutes()
	//Serve assets

	//router.Add("/login", controller.LoginHandler)

	handler := &router.Handler{}
	fmt.Println("Starting server on :8080")
	if err := http.Serve(listener, handler); err != nil {
		panic(err)
	}
}
