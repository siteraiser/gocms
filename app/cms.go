package app

import (
	"bytes"
	"fmt"
	"gocms/app/helpers"
	"net/http"
	"net/url"
)

/*
Main interface using app.Cms{}
Uses some app.functions to actually get the values.
*/
// -------------------

func Cms(r *http.Request) cms {
	return cms{
		Header: Header{R: r},
		URL:    URL{R: r},
		Any: URLAnyValue{
			R:    r,
			Vals: AnyValues(r),
		},
		Named: URLNameValue{R: r},
		Form:  Form{R: r},
		Views: Views{R: r},
		Utils: Utils{},
	}
}

type cms struct {
	Header Header
	URL    URL
	Any    URLAnyValue
	Named  URLNameValue
	Form   Form
	Views  Views
	Utils  Utils
}

type Header struct {
	R *http.Request
}

type URL struct {
	R *http.Request
}
type QueryParams struct {
	R *http.Request
}

type URLAnyValue struct {
	Vals []string
	R    *http.Request
}
type URLNameValue struct {
	Vals map[string]string
	R    *http.Request
}

type Form struct {
	R *http.Request
}

type Views struct {
	R *http.Request
}

type Utils struct {
	Html helpers.Html
}

func (r URL) Path() string {
	return r.R.URL.Path
}

// Header Presets --------
// Add custom header
func (r Header) Set(key string, value string) cms {
	Requests[GetId(r.R)].Headers = append(Requests[GetId(r.R)].Headers, [2]string{key, value})
	return Cms(r.R)
}

type ResponseTypesList struct {
	TextPlain [2]string
	TextHtml  [2]string
	TextJson  [2]string
}

/* Different preset response type headers*/
var ResponseTypes = ResponseTypesList{
	TextPlain: [2]string{"Content-Type", "text/plain; charset=utf-8"},
	TextHtml:  [2]string{"Content-Type", "text/html; charset=utf-8"},
	TextJson:  [2]string{"Content-Type", "text/json; charset=utf-8"},
}

// Set response header to: "Content-Type", "text/plain; charset=utf-8"
func (r Header) SetTextPlain() cms {
	setType(r, ResponseTypes.TextPlain, Requests[GetId(r.R)].Headers)
	return Cms(r.R)
}

// Set response header to: "Content-Type", "text/html; charset=utf-8"
func (r Header) SetTextHtml() cms {
	setType(r, ResponseTypes.TextHtml, Requests[GetId(r.R)].Headers)
	return Cms(r.R)
}

// Set response header to: "Content-Type", "text/json; charset=utf-8"
func (r Header) SetTextJson() cms {
	setType(r, ResponseTypes.TextJson, Requests[GetId(r.R)].Headers)
	return Cms(r.R)
}

// helper for header set methods (looks for header and replaces if need be)
func setType(r Header, typetuple [2]string, headers [][2]string) {
	for i, header := range headers {
		if header[0] == "Content-Type" {
			Requests[GetId(r.R)].Headers[i] = typetuple
			return
		}
	}
	Requests[GetId(r.R)].Headers = append(Requests[GetId(r.R)].Headers, typetuple)
}

// Add presets maybe eg, html page, json api etc
// Request Parameters --------
// Returns url.Values from "url.Parse" function
func (r URL) QueryParams() url.Values {
	parsedURL, err := url.Parse("?" + r.R.URL.RawQuery)
	if err != nil {
		fmt.Println("Error:", err)
		return url.Values{}
	}
	return parsedURL.Query()
}

// Returns a query parameter value by key
func (r URL) QueryParam(i string) string {
	parsedURL, err := url.Parse("?" + r.R.URL.RawQuery)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	queryParams := parsedURL.Query()
	return queryParams.Get(i)
}

// Returns an array of the URL segments
func (r URL) Segments() []string {
	return UrlSegments(r.R)
}

// Returns any value "{$}" by index
func (r URLAnyValue) Value(i int) string {
	return AnyValue(r.R, i)
}

// Returns array of any values "{$}"
func (r URLAnyValue) Values() []string {
	return AnyValues(r.R)
}

// Returns map of named values "{id}"
func (r URLNameValue) Value(i string) string {
	return NameValue(r.R, i)
}

// Returns anamedny value "{id}" by string key
func (r URLNameValue) Values() map[string]string {
	return NamedValues(r.R)
}

// Returns form fields from post request
func (r Form) Fields() url.Values {
	err := r.R.ParseForm()
	if err != nil {
		fmt.Println("Error:", err)
		return url.Values{}
	}
	return r.R.Form
}

// Returns form field from post request by key
func (r Form) Field(i string) string {
	err := r.R.ParseForm()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return r.R.Form.Get(i)
}

// Views --------
// Calls render function
func (v Views) Render(location string, args any) string {
	return Render(location, args)
}

// Calls render function, appends to Response Ouptut and returns the string result
func (v Views) Add(location string, args any) string {
	content := Render(location, args)
	contentbytes := []byte(content)
	Requests[GetId(v.R)].Output = append(Requests[GetId(v.R)].Output, contentbytes...)
	return content
}

// Returns the view output from all prior Views.Add  calls
func (v Views) OutBytes() []byte {
	return Requests[GetId(v.R)].Output
}
func (v Views) Out() string {
	buffer := bytes.NewBuffer(Requests[GetId(v.R)].Output)
	return buffer.String()
}

// Set output (overwrites content)
func (v Views) SetOut(content string) string {
	contentbytes := []byte(content)
	Requests[GetId(v.R)].Output = contentbytes
	return content
}

// Append to output
func (v Views) OutAppend(content string) string {
	contentbytes := []byte(content)
	Requests[GetId(v.R)].Output = append(Requests[GetId(v.R)].Output, contentbytes...)
	return content
}

// Utils --------
// Append to output
func (h Utils) Ahref(href string, text string) string {
	return helpers.Html.Ahref(helpers.Html{}, href, text)
}
