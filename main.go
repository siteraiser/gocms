package main

import (
	"fmt"
	"gocms/app"
	"gocms/app/router"
	"net"
	"net/http"
)

func main() {

	//portFlag := flag.Int("port", 8080, "string")
	//flag.Parse()
	port := "8080" //strconv.Itoa(*portFlag)
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		panic(err)
	}

	app.NewApp(&router.Handler{})

	//Add routes, serve assets
	addRoutes()

	fmt.Println("Starting server on :" + port)
	if err := http.Serve(listener, app.Router.Handler); err != nil {
		panic(err)
	}
}
