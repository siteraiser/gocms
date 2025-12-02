package app

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	db "gocms/app/database"
	"gocms/app/models"
	"gocms/app/router"
	"gocms/app/sys"
	"slices"

	"gocms/templates"
	"net/http"
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
var Config sys.Configuration
var Db *sql.DB

func NewApp(h http.Handler) {
	Router = Routing{
		h,
	}

	Config = sys.Config
	BaseUrl = Config.Settings.BaseUrl
	//Start stats display... maybe add some options to config
	if Config.Settings.StatsEnabled {
		sys.SysStats()
		sys.Stats.ReqRef = Requests
	}
	//Start Db
	if Config.Database.UserName != "" {
		Db = db.InitDB(db.DbConfig{
			UserName: Config.Database.UserName,
			Password: Config.Database.Password,
			Host:     Config.Database.Host,
			Port:     Config.Database.Port,
			DbName:   Config.Database.Name,
		})
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
	err := GetPage(w, r.Clone(ctx))
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
		} else {
			route = router.Route{}
		}
	}

	//if still not found then check for auto-routed controllers/functions (package/function)...
	var hfn http.HandlerFunc
	if !found && Config.Settings.AutoRoutes {
		hfn, found = router.AutoRouteIt(path, urlsegs)
		if found {
			routetype = "auto"
		} else {
			hfn = nil
		}
	}

	if found {
		//	Mutex.Lock()
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
		//	Mutex.Unlock()

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
					if slices.Contains(strings.Split(r.Header.Get("accept-encoding"), ","), "gzip") && len(*output) != 0 {
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

		fmt.Printf("\nServed from %v router:  %v in %v\n", routetype, path, time.Since(start))
		//fmt.Println("\nRoute: ", route)
		//If combining content from multiple views, flush after serving

		flusher.Flush()
		// Do background work without blocking the client
		go func() {
			//consider r lock
			delete(Requests, requestid)
			Mutex.Lock()
			sys.Stats.TotalHits++
			Mutex.Unlock()

		}()
		return
	}
	fmt.Println("Not found with router: ", path)
	//	fmt.Println("\nRoute: ", route)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Custom 404: Page not found"))
	flusher.Flush()
	go func() {
		delete(Requests, requestid)
		Mutex.Lock() //consider r lock
		sys.Stats.TotalHits++
		Mutex.Unlock()

	}()
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

var Requests = map[string]*models.Request{}

// ------------------------------
// Common functions via app.*
// ------------------------------
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
