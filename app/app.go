package app

import (
	"bytes"
	"context"
	"fmt"
	"slices"

	"gocms/app/models"
	"gocms/app/router"
	"gocms/templates"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"compress/gzip"

	"github.com/google/uuid"
)

type key string

const RequestIDKey key = "requestID"

var Mutex sync.Mutex

// add routing from routing package to app
type Routing struct {
	Handler http.Handler
}

var Router Routing

func NewApp(h http.Handler) {
	Router = Routing{
		h,
	}
}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestid := r.Header.Get("X-Request-Id")
	if requestid == "" {
		requestid = uuid.New().String()
	}

	ctx := context.WithValue(r.Context(), "RID", requestid)
	w.Header().Set("X-Request-Id", requestid)
	//w.Header().Add("requestid", requestid)

	req := models.Request{
		Id:      requestid,
		Handler: h,
		Path:    r.URL.Path,
	}
	Mutex.Lock()
	Requests[requestid] = &req
	Mutex.Unlock()

	path := r.URL.Path
	if path != "/" {
		path = strings.TrimLeft(path, "/")
	}
	urlsegs := strings.Split(path, "/")

	routetype := ""
	found := false

	//Capture the output and send it but clear the output on the way out
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	//check for other resources that aren't using the default routing here
	req.Path = path
	req.UrlSegments = urlsegs

	//try custom / primary...
	err := router.GetPage(w, r.Clone(ctx))
	if err == nil {
		routetype = "primary"
		found = true
	}

	//try other ones...
	//No manual checks matched the request URL using the primary router, now try secondary router
	var route router.Route
	var anyvalues []string
	var namedvalues map[string]string
	if !found {
		route, anyvalues, namedvalues, found = router.RouteIt(path, r.Method)
		if found {
			routetype = "secondary"
		}
	}

	//if still not found then check for auto-routed controllers/functions (package/function)...
	var hfn http.HandlerFunc
	if !found && Config.Settings.AutoRoutes {
		hfn, found = router.AutoRouteIt(path, urlsegs)
		if found {
			routetype = "auto"
		}
	}

	if found {
		Mutex.Lock()
		Requests[requestid].Path = path
		Requests[requestid].AnyValues = anyvalues
		Requests[requestid].NamedValues = namedvalues
		Requests[requestid].RouteType = routetype

		if routetype == "secondary" {
			Requests[requestid].Handler = route.Controller
		}
		if routetype == "auto" {
			Requests[requestid].HandlerFunc = hfn
		}
		Mutex.Unlock()

		for _, req := range Requests {
			if req.Id == requestid {
				switch req.RouteType {
				case "secondary":
					req.Handler.ServeHTTP(w, r.WithContext(ctx))
				case "auto":
					req.HandlerFunc.ServeHTTP(w, r.WithContext(ctx))
				}

				if Config.Settings.UseViewOutput == true {

					if len(Requests[requestid].Headers) == 0 {
						w.Header().Set("Content-Type", "text/html; charset=utf-8")
					} else {
						for _, header := range Requests[requestid].Headers {
							w.Header().Set(header[0], header[1])
						}
					}

					output := &Requests[requestid].Output
					if slices.Contains(strings.Split(r.Header.Get("accept-encoding"), ","), "gzip") && routetype != "primary" && len(*output) != 0 {
						// Set headers for gzip encoding
						w.Header().Set("Content-Encoding", "gzip")
						w.Header().Del("Content-Length") // Remove Content-Length as it changes after compression
						*output = gzipper(output)
					}
					_, err = w.Write(*output)
					Requests[requestid].Output = nil
					if err != nil {
						fmt.Println("Error writing response:", err)
					}
				}

			}

		}

		fmt.Printf("Served from %v router:  %v in %v\n", routetype, path, time.Since(start))
		fmt.Println("route: ", route)
		//If combining content from multiple views, flush after serving

		flusher.Flush()
		// Do background work without blocking the client
		go func() {
			//ClearOutput(requestid)
			Mutex.Lock() //consider r lock
			delete(Requests, requestid)
			Mutex.Unlock()
		}()
		return
	}
	fmt.Println("Not found with router: ", path)
	fmt.Println("route: ", route)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Custom 404: Page not found"))
	flusher.Flush()

}

func gzipper(a *[]byte) []byte {

	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	defer gz.Close()
	if _, err := gz.Write(*a); err != nil {
		panic(err)
	}
	gz.Close()
	//	gz.Flush()
	return b.Bytes()
}

var BaseUrl = ""

var Requests = make(map[string]*models.Request)

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
	}
}

type cms struct {
	Header Header
	URL    URL
	Any    URLAnyValue
	Named  URLNameValue
	Form   Form
	Views  Views
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

func (r URL) Path() string {
	return r.R.URL.Path
}

// Set header type for correct output
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

// "Content-Type", "text/plain; charset=utf-8"
func (r Header) SetTextPlain() cms {
	setType(r, ResponseTypes.TextPlain, Requests[GetId(r.R)].Headers)
	return Cms(r.R)
}

// "Content-Type", "text/html; charset=utf-8"
func (r Header) SetTextHtml() cms {
	setType(r, ResponseTypes.TextHtml, Requests[GetId(r.R)].Headers)
	return Cms(r.R)
}

// "Content-Type", "text/json; charset=utf-8"
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
func (v Views) Out() []byte {
	return Requests[GetId(v.R)].Output
}

// ---
func AnyValue(r *http.Request, index int) string {
	Any := AnyValues(r)
	if len(Any)-1 >= index {
		return Any[index]
	}
	return ""
}
func NameValue(r *http.Request, index string) string {
	Named := NamedValues(r)
	if vals, exists := Named[index]; exists {
		return vals
	}
	return ""
}

// ------------------------------

// ------------------------------
// Shortcuts to common functions via app.*

func GetId(r *http.Request) string {
	//fmt.Println("r.Context().Value():", r.Context().Value("RID"))
	if r.Context().Value("RID") != nil {
		return r.Context().Value("RID").(string)
	}
	return ""
}

func UrlSegments(r *http.Request) []string {
	return Requests[GetId(r)].UrlSegments
}

func AnyValues(r *http.Request) []string {
	return Requests[GetId(r)].AnyValues
}

func NamedValues(r *http.Request) map[string]string {
	return Requests[GetId(r)].NamedValues
}

func Render(location string, args any) string {
	out := ""
	//no reason to choose engine for now with: app.GetConfig()...
	data, err := os.ReadFile("./views/" + location)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	err = nil
	//Find the rendering engine in the registry (outside of app folder) and render
	for _, e := range templates.Registry {
		if e.Ext == filepath.Ext(location) {
			out, err = e.Engine.Render(string(data), args)

			//	Requests[GetId(w)].View.Output = out

			if err != nil {
				fmt.Println("Error:", err)
				return ""
			}
		}
	}

	//add more types of rendering here...
	return out
}

// --- other random helpers...
func Ahref(href string, text string) string {
	return "<a href='" + href + "'>" + text + "</a>"
}
