package models

import "net/http"

type Request struct {
	Id          string
	RouteType   string
	Handler     http.Handler
	HandlerFunc http.HandlerFunc
	Path        string
	UrlSegments []string
	AnyValues   []string
	NamedValues map[string]string
}
