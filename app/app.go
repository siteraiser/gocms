package app

import (
	"net/http"
)

type App struct {
	Router Router
}

type Router struct {
	Handler http.Handler
}

func NewApp(ah http.Handler) App {
	return App{Router{ah}}
}

type Page struct {
	Attributes struct {
		Header http.Header
	}
	Meta    string
	Content string
	Assets  map[string]string
}
