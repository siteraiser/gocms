package models

import "net/http"

type Request struct {
	Id           string //UUID for x-request header
	RouteType    string
	ResponseType [2]string
	Handler      http.Handler
	HandlerFunc  http.HandlerFunc
	Path         string
	UrlSegments  []string
	AnyValues    []string
	NamedValues  map[string]string
	Output       string
}
