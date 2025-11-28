package app

import (
	"context"
	"fmt"

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
	if !found && Config.Settings.Preferences.AutoRoutes {
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

var BaseUrl = ""

var Requests = make(map[string]*models.Request)

type URLValues struct {
	AnyValues   []string
	NamedValues map[string]string
}

type RequestValues struct {
	AnyValue   func(int) string
	NamedValue func(string) string
	URLVals    URLValues
}

/*
Alternate interface using app.Cms{}
Uses app.functions to actually get the values.

type URLInterface interface {
	AnyValues(i int) []string
	NamedValues(i string) map[string]string
}

*/
// -------------------

func Cms(r *http.Request) cms {
	return cms{
		URL:   URL{R: r},
		Any:   URLAnyValue{R: r},
		Named: URLNameValue{R: r},
		Form:  Form{R: r},
		Views: Views{R: r},
	}
}

type cms struct {
	URL   URL
	Any   URLAnyValue
	Named URLNameValue
	Form  Form
	Views Views
}

type URL struct {
	R *http.Request
}
type QueryParams struct {
	R *http.Request
}

type URLAnyValue struct {
	R *http.Request
}
type URLNameValue struct {
	R *http.Request
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

func (r URL) QueryParams() url.Values {
	parsedURL, err := url.Parse("?" + r.R.URL.RawQuery)
	if err != nil {
		fmt.Println("Error:", err)
		return url.Values{}
	}
	return parsedURL.Query()
}
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

// Calls render function saves the ouptut and returns the string result
func (v Views) Add(location string, args any) string {
	content := AddView(location, args)
	Requests[GetId(v.R)].Output = content
	return content
}

// Returns the view output from all prior Add view calls
func (v Views) Out() string {
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
func Req(r *http.Request) URLValues {
	return URLValues{
		AnyValues:   AnyValues(r),
		NamedValues: NamedValues(r),
	}
}

func GetId(r *http.Request) string {
	fmt.Println("r.Context().Value():", r.Context().Value("RID"))
	if r.Context().Value("RID") != nil {
		return r.Context().Value("RID").(string)
	}
	return ""
}

func UrlSegments(r *http.Request) []string {
	fmt.Println("GetId(r):", GetId(r))
	return Requests[GetId(r)].UrlSegments
}

func AnyValues(r *http.Request) []string {
	fmt.Println("AnyValues:", Requests)
	return Requests[GetId(r)].AnyValues
}

func NamedValues(r *http.Request) map[string]string {
	fmt.Println("NamedValues:", Requests)
	return Requests[GetId(r)].NamedValues
}

func AddView(location string, args any) string {
	return Render(location, args)
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
