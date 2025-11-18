package router

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct{}

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

type Route struct {
	Pattern    string
	Controller http.Handler
}
type Routes struct {
	List []Route
}

var routes Routes
var AnyParams []string
var NamedParams map[string]string

func AnyValues() []string {
	return AnyParams
}
func AnyValue(index any) string {
	if index != nil {
		if i, err := strconv.Atoi(index.(string)); err == nil {
			return AnyParams[i]
		}
	}
	return ""
}
func NamedValue(name string) string {
	return NamedParams[name]
}

func Add(pattern string, controller http.Handler) {
	var route Route
	route.Pattern = pattern
	route.Controller = controller
	routes.List = append(routes.List, route)
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check for other resources that aren't using the default routing here
	err := GetPage(w, r)
	if err == nil {
		return
	}
	//No manual checks matched the request URL, now try router
	path := strings.TrimLeft(r.URL.Path, "/")

	route, anyParams, namedParams, found := routeIt(path)
	AnyParams = anyParams
	NamedParams = namedParams

	if found {
		fmt.Println("Served from router: ", r.URL.Path)
		route.Controller.ServeHTTP(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Custom 404: Page not found"))
}

func GetPage(w http.ResponseWriter, r *http.Request) error {
	//Overrides normal app routes
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div>HOME:%v\n <a href='/task/1'>task 1</a> <a href='/app/1'>app 1</a></div>", r.URL.Path)

	} else if r.URL.Path == "/test" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, "<div><img src='/assets/media/images/pic.png'>test ok path:%v\n <a href='/'>home</a> <a href='/task/1'>task 1</a></div>", r.URL.Path)
	} else {
		err := errors.New("Something went wrong")
		return err
	}
	fmt.Println("Served from primary routes: ", r.URL.Path)
	return nil

}

func routeIt(path string) (Route, []string, map[string]string, bool) {

	match := func(pattern []string, url_segs []string) ([]string, map[string]string, bool) {
		i := 0
		match := false
		urlcount := len(url_segs)
		anything := []string{}
		named := map[string]string{}
		for _, value := range pattern {
			pattern_param := string(value)
			if len(url_segs) > i && len(pattern_param) > 0 {

				if pattern_param[0:1] != "{" && pattern_param[len(pattern_param)-2:len(pattern_param)-1] != "}" {

					if url_segs[i] == pattern_param {
						match = true
					} else {
						return []string{}, map[string]string{}, false
					}
				} else {
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
	var found bool
	var route Route
	//var routeNo int
	for _, route = range routes.List {

		url_segs := strings.Split(path, "/")
		pattern := strings.Split(strings.TrimLeft(string(route.Pattern), "/"), "/")

		p_len := len(pattern)
		if p_len > len(url_segs) {
			p_len = len(url_segs)
		}

		//Step through the pattern and the path simultaneously to look for matches
		anyParams, namedParams, found = match(pattern, url_segs[:p_len])
		// Should be ok lol, maybe further investigate
		if string(route.Pattern)[len(string(route.Pattern))-1:len(string(route.Pattern))] == "/" {
			found = true
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
