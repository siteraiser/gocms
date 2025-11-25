package router

import (
	"errors"
	"fmt"
	"gocms/app"
	"net/http"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

/*
	type Response struct {
		Page Page
	}

	type Page struct {
		Attributes struct {
			Header http.Header
		}
		Meta    string
		Content string
		Assets  map[string]string
	}
*/
type Route struct {
	Pattern     string
	Controller  http.Handler
	RequestType string
}
type Routes struct {
	List []Route
}

var routes Routes

type Controllers struct {
	List map[string]http.HandlerFunc
}

var controllers Controllers

// Add a route
func Add(pattern string, controller http.Handler) {
	var route Route

	route.Controller = controller
	if strings.HasPrefix(pattern, "GET ") {
		route.RequestType = "GET"
		route.Pattern = pattern[4:]
	} else if strings.HasPrefix(pattern, "POST ") {
		route.RequestType = "POST"
		route.Pattern = pattern[5:]
	} else {
		route.RequestType = ""
		route.Pattern = pattern
	}

	routes.List = append(routes.List, route)
}

// Add a controller
func AddFunc(controller http.HandlerFunc) {
	funcValue := reflect.ValueOf(controller)
	fmt.Println("Value of controller: ", funcValue)
	// simple mvc routing
	name := getName(controller)
	fmt.Println("name of controller: ", name)
	if controllers.List == nil {
		controllers.List = make(map[string]http.HandlerFunc)
	}
	if name != "" {
		before, fn, _ := strings.Cut(name, ".")
		pkgArr := strings.Split(before, "/")
		pkg := pkgArr[len(pkgArr)-1]
		fmt.Println("pkg", pkg)
		fmt.Println("function", fn)
		controllers.List[strings.ToLower(pkg+"/"+fn)] = controller
	}

}

func getName(myvar interface{}) string {

	return runtime.FuncForPC(reflect.ValueOf(myvar).Pointer()).Name()

}

type Handler struct{}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	if path != "/" {
		path = strings.TrimLeft(path, "/")
	}
	urlsegs := strings.Split(path, "/")
	//app.Request = r

	//Capture the output and send it but clear the output on the way out
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}
	routeType := ""
	found := false
	//check for other resources that aren't using the default routing here
	err := GetPage(w, r)
	if err == nil {
		routeType = "primary"
		found = true
	}

	//No manual checks matched the request URL, now try router
	var route Route
	var anyvalues []string
	var namedvalues map[string]string
	if !found {
		route, anyvalues, namedvalues, found = routeIt(path, r.Method)
		if found {
			routeType = "secondary"
		}
	}

	//if still not found then check for auto-routed controllers/functions (package/function)...
	var hfn http.HandlerFunc
	if !found && app.Config.Settings.Preferences.AutoRoutes {
		hfn, found = autoRouteIt(path, urlsegs)
		if found {
			routeType = "auto"
		}
	}

	if found {
		fmt.Println("Served from "+routeType+" router: ", app.Path)
		fmt.Println("route: ", route)
		app.Mutex.Lock()
		app.AnyValues = anyvalues
		app.NamedValues = namedvalues
		if routeType == "secondary" {
			route.Controller.ServeHTTP(w, r)
		}
		if routeType == "auto" {
			hfn.ServeHTTP(w, r)
		}
		app.Mutex.Unlock()

		//If combining content from multiple views, flush after serving
		flusher.Flush()

		// Do background work without blocking the client
		go func() {
			app.ClearOutput()

		}()
		return
	}

	fmt.Println("Not found with router: ", path)
	fmt.Println("route: ", route)

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Custom 404: Page not found"))
	flusher.Flush()

}

func GetPage(w http.ResponseWriter, r *http.Request) error {
	//Overrides normal app routes
	/*if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>HOME:%v\n <a href='/task/1'>task 1</a> <a href='/app/1'>app 1</a></div>", r.URL.Path)

	} else
	*/
	if r.URL.Path == "/test" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div><img src='/assets/media/images/pic.png'>test ok path:%v\n <a href='/'>home</a> <a href='/another/value1/and/value2'>Test page</a></div>", r.URL.Path)
	} else {
		err := errors.New("Something went wrong")
		return err
	}
	fmt.Println("Served from primary routes: ", r.URL.Path)
	return nil

}

func routeIt(path string, method string) (Route, []string, map[string]string, bool) {

	match := func(pattern []string, url_segs []string) ([]string, map[string]string, bool) {
		i := 0
		match := false
		urlcount := len(url_segs)
		anything := []string{}
		named := map[string]string{}
		for _, value := range pattern {
			pattern_param := string(value)
			if len(url_segs) > i && len(pattern_param) > 0 {
				//if segment is not a parameter
				if pattern_param[0:1] != "{" && pattern_param[len(pattern_param)-2:len(pattern_param)-1] != "}" {
					//if it matches
					if url_segs[i] == pattern_param {
						match = true
					} else {
						return []string{}, map[string]string{}, false
					}
				} else {
					//Save any {$} and named {id} parameter values
					if url_segs[i] != "" {
						parametervalue := pattern_param[1 : len(pattern_param)-1]
						if parametervalue == "$" {
							anything = append(anything, url_segs[i])
						} else if len(parametervalue) > 0 {
							named[parametervalue] = url_segs[i]
						}
					}
				}
			} else {
				return []string{}, map[string]string{}, false
			}
			i++

			if urlcount == 0 {
				return []string{}, map[string]string{}, false
			}
			urlcount--
		}
		if len(anything) == 0 && len(named) == 0 && match {
			return []string{}, map[string]string{}, true
		}
		return anything, named, true
	}

	i := 0
	var anyParams []string
	var namedParams map[string]string
	found := false
	var route Route

	for _, route = range routes.List {

		if path == "/" && route.Pattern == "/" {
			//index
			found = true
			break
		}
		url_segs := strings.Split(path, "/")
		pattern := strings.Split(strings.TrimLeft(string(route.Pattern), "/"), "/")

		pattern_str := string(route.Pattern)
		pattern_str_len := len(string(route.Pattern))

		//Filter by request type
		if route.RequestType != "" && route.RequestType != method {
			found = false
			break
		}

		//Check if the pattern is a "folder/" and request is a file
		if pattern_str[pattern_str_len-1:pattern_str_len] == "/" &&
			strings.Contains(filepath.Base(path), ".") { //should be improved (maybe add allowed file types list...)
			if "/"+path[0:pattern_str_len-1] != pattern_str {
				// If is a file...
				found = false
			} else {
				found = true
			}
			break
		} else {

			p_len := min(len(pattern), len(url_segs))
			//Step through the pattern and the path simultaneously to look for matches
			anyParams, namedParams, found = match(pattern, url_segs[:p_len])
		}

		if found {
			break
		}
		i++
	}

	var anys []string
	var named map[string]string
	if len(anyParams) > 0 {
		anys = append(anys, anyParams...)
	}
	if len(namedParams) > 0 {
		named = namedParams
	}

	return route, anys, named, found
}

func autoRouteIt(path string, urlsegs []string) (http.HandlerFunc, bool) {
	//found := false
	controller_name := ""
	package_name := ""
	if path == "/" {
		controller_name = "index"
	} else if len(urlsegs) > 1 {
		controller_name = urlsegs[1]
		package_name = urlsegs[0]
	}
	if _, exists := controllers.List[package_name+"/"+controller_name]; !exists {
		controller_name = "index"
		package_name = urlsegs[0]
	}
	if _, exists := controllers.List[package_name+"/"+controller_name]; !exists {
		controller_name = "index"
		package_name = "controller"
	}
	mvcroute := package_name + "/" + controller_name
	fmt.Println("mvcroute: ", mvcroute)
	if fn, exists := controllers.List[mvcroute]; exists {

		return fn, true
	} else {
		fmt.Println("Package not found", controller_name)
	}
	//return found
	return nil, false
}
