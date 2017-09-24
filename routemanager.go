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
	parts := str.Split(request.URL.EscapedPath(), '/')

	controllerName := manager.DefaultController
	actionName := manager.DefaultAction
	params := []string{}

	if len(parts) > 2 {
		controllerName = parts[0]
		actionName = parts[1]
		params = parts[2:]
	} else {
		if len(parts) >= 2 {
			controllerName = parts[0]
			actionName = parts[1]
		} else {
			if len(parts) == 1 {
				controllerName = parts[0]
			}
		}
	}

	for _, route := range manager.Routes {
		if str.Compare(controllerName, route.ControllerName) {
			// Construct the appropriate controller
			icontroller := route.CreateController(request)
			controller := icontroller.ToController()

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
			result := icontroller.Execute(actionName, params)

			// Write controllers cookies
			for _, cookie := range controller.Cookies {
				http.SetCookie(response, cookie)
			}

			// Execute the response and return
			// TODO: Handle Errors here
			result.Execute(response)
			return
		}
	}

	// TODO:
	// Handle 404 (No route found)
}
