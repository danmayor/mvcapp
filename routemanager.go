/*
	Digivance MVC Application Framework
	Route Manager Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the generic route manager functionality of the MVC application. This manager
	allows the caller to register route maps and bind the handler method. This system drives the
	request pipeline of an MVC application made with this framework.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Digivance/str"
)

// ControllerCreator is a delegate to the creation method of a controller
// that is mapped as the primary route. (E.g. site.com/CONTROLLER is mapped
// to a NewXController method that implements this signature)
type ControllerCreator func(*http.Request) IController

// RouteManager provides the basic http request pipeline of the
// mvcapp framework
type RouteManager struct {
	// SessionIDKey is the name of the cookie used to define the browser session ID
	// for incoming requests
	SessionIDKey string

	DefaultController string
	DefaultAction     string

	Routes         []*RouteMap
	SessionManager *SessionManager
}

// NewRouteManager returns a new route manager object with default
// controller and action tokens set to "Home" and "Index".
func NewRouteManager() *RouteManager {
	return &RouteManager{
		SessionIDKey: "MvcApp.SessionID",

		DefaultController: "Home",
		DefaultAction:     "Index",

		Routes: make([]*RouteMap, 0),
	}
}

// RegisterController is used to map a custom controller object to the
// controller section of the requested url (E.g. "site.com/CONTROLLER/action")
func (manager *RouteManager) RegisterController(name string, creator ControllerCreator) {
	manager.Routes = append(manager.Routes, NewRouteMap(name, creator))
}

// HandleRequest is mapped to the http handler method and processes the
// HTTP request pipeline
func (manager *RouteManager) HandleRequest(response http.ResponseWriter, request *http.Request) {
	if manager.handleController(response, request) {
		return
	}

	manager.handleFile(response, request)
}

// parseFragment will extract, and remove the fragment (or named anchor) section
// of the url. Returns strings representing the fragment and the url without it
func (manager *RouteManager) parseFragment(url string) (string, string) {
	fragment := ""

	if str.Contains(url, "#") {
		fragment = str.RightOf(url, "#")
		url = str.LeftOf(url, "#")
	}

	return fragment, url
}

// parseQueryString will extract the path and parse the query string key value pairs.
// Returns the path (relative to the app's domain) and a map of the query string pairs.
func (manager *RouteManager) parseQueryString(url string) (string, map[string]string) {
	path := ""
	queryString := map[string]string{}

	if str.Contains(url, "?") {
		path = str.LeftOf(url, "?")
		qsLine := str.RightOf(url, "?")

		for _, pair := range str.Split(qsLine, '&') {
			kvp := str.Split(pair, '=')
			if len(kvp) > 1 {
				queryString[kvp[0]] = kvp[1]
			}

			if len(kvp) == 1 {
				queryString[kvp[0]] = ""
			}
		}
	} else {
		path = url
	}

	return path, queryString
}

// parseControllerName returns the controller name requested, will fallback and return
// the default controller if this is a root request.
func (manager *RouteManager) parseControllerName(path string) string {
	rtn := manager.DefaultController
	parts := str.Split(path, '/')

	if len(parts) > 0 {
		rtn = parts[0]
	}

	return rtn
}

func (manager *RouteManager) handleController(response http.ResponseWriter, request *http.Request) bool {
	fragment, url := manager.parseFragment(request.URL.Path)
	path, queryString := manager.parseQueryString(url)
	controllerName := manager.parseControllerName(path)

	for _, route := range manager.Routes {
		if str.StartsWith(route.ControllerName, controllerName) {
			// Construct the appropriate controller
			icontroller := route.CreateController(request)
			controller := icontroller.ToController()

			controller.DefaultAction = manager.DefaultAction
			controller.RequestedPath = path
			controller.QueryString = queryString
			controller.Fragment = fragment

			if manager.SessionManager != nil {
				// Get the browser session ID from the request cookies
				browserSessionCookie, err := request.Cookie(manager.SessionIDKey)
				browserSessionID := ""
				if err != nil {
					browserSessionID = str.Random(32)
				} else {
					browserSessionID = browserSessionCookie.Value
				}

				// Get the browserSession from the SessionManager and set
				// the controllers reference to it
				browserSession := manager.SessionManager.GetSession(browserSessionID)
				controller.Session = browserSession
				controller.Session.ActivityDate = time.Now().Add(900 * time.Second)
			}

			// Prepare result
			if controller.BeforeExecute != nil {
				controller.BeforeExecute()
			}

			result := icontroller.Execute()

			if controller.AfterExecute != nil {
				controller.AfterExecute()
			}

			// Write controllers cookies
			for _, cookie := range controller.Cookies {
				http.SetCookie(response, cookie)
			}

			// Execute the response and return
			result.Execute(response)
			return true
		}
	}

	return false
}

func (manager *RouteManager) handleFile(response http.ResponseWriter, request *http.Request) {
	_, url := manager.parseFragment(request.URL.Path)
	path, _ := manager.parseQueryString(url)

	if str.StartsWith(path, "/") {
		path = fmt.Sprintf("%s/%s", GetApplicationPath(), path[1:])
	}

	if validPath(path) {
		http.ServeFile(response, request, path)
	}
}

func validPath(path string) bool {
	if str.StartsWith(path, "controllers/") {
		return false
	}

	if str.StartsWith(path, "models/") {
		return false
	}

	if str.StartsWith(path, "views/") {
		return false
	}

	if !str.Contains(path, "/") {
		return false
	}

	return true
}
