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
	"os"
	"time"

	"github.com/Digivance/applog"
	"github.com/Digivance/str"
)

// ControllerCreator is a delegate to the creation method of a controller
// that is mapped as the primary route. (E.g. site.com/CONTROLLER is mapped
// to a NewXController method that implements this signature)
type ControllerCreator func(*http.Request) IController

// RouteManager provides the basic http request pipeline of the
// mvcapp framework
type RouteManager struct {
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

// handleController is used to attempt to handle a request through an mvcapp controller
// pipeline. (E.g. controller.BeforeExecute, controller.Execute, writes the header with
// the controller.HTTPStatusCode and then executes the ActionResult (to deliver the payload)
// and finally the controller.AfterExecute
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

			// Write controllers cookies
			for _, cookie := range controller.Cookies {
				http.SetCookie(response, cookie)
			}

			// Prepare result
			if controller.BeforeExecute != nil {
				controller.BeforeExecute()
			}

			result := icontroller.Execute()

			if result == nil {
				if controller.NotFoundResult != nil {
					result = controller.NotFoundResult(controller.RequestedPath)
				} else {
					result = controller.DefaultNotFoundPage()
				}
			}

			if err := result.Execute(response); err != nil {
				msg := err.Error()
				if str.Compare(msg, "No response from request") {
					if controller.NotFoundResult != nil {
						result = controller.NotFoundResult(controller.RequestedPath)
					} else {
						result = controller.DefaultNotFoundPage()
					}

					if err = result.Execute(response); err != nil {
						applog.WriteError("Failed to display default 404 page!", err)
						response.WriteHeader(404)
					}
				} else {
					if controller.ErrorResult != nil {
						result = controller.ErrorResult(err)
					} else {
						result = controller.DefaultErrorPage(err)
					}

					if err = result.Execute(response); err != nil {
						applog.WriteError("Failed to display default error page!", err)
						response.WriteHeader(500)
					}
				}
			}

			if controller.AfterExecute != nil {
				controller.AfterExecute()
			}

			return true
		}
	}

	return false
}

// handleFile is called if handleController is false and will attempt to serve a raw file
func (manager *RouteManager) handleFile(response http.ResponseWriter, request *http.Request) {
	_, url := manager.parseFragment(request.URL.Path)
	path, _ := manager.parseQueryString(url)

	if str.StartsWith(path, "/") {
		path = fmt.Sprintf("%s/%s", GetApplicationPath(), path[1:])
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		for _, route := range manager.Routes {
			if str.StartsWith(route.ControllerName, manager.DefaultController) {
				controller := route.CreateController(request).ToController()
				if controller.NotFoundResult != nil {
					result := controller.NotFoundResult(url)
					if err := result.Execute(response); err != nil {
						response.WriteHeader(404)
					}
				} else {
					result := controller.DefaultNotFoundPage()
					if err := result.Execute(response); err != nil {
						response.WriteHeader(404)
					}
				}
			}
		}
	} else {
		if validPath(path) {
			http.ServeFile(response, request, path)
		}
	}
}

// validPath is used internally to ignore paths that are used by the mvcapp system
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
