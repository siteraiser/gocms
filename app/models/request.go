package models

import "net/http"

type Request struct {
	Id           string //UUID for x-request header
	RouteType    string
	ResponseType string
	Headers      [][2]string
	Handler      http.Handler
	HandlerFunc  http.HandlerFunc
	Path         string
	UrlSegments  []string
	AnyValues    []string
	NamedValues  map[string]string
	Output       string
}

type ResponseTypesList struct {
	TextPlain [2]string
	TextHtml  [2]string
	TextJson  [2]string
}
