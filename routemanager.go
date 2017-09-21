/*
	Digivance MVC Application Framework
	Route Manager Features
	Dan Mayor (dmayor@digivance.com)

	This file defines the basic route manager functionality
*/

package mvcapp

import (
	"net/http"

	"github.com/Digivance/str"
)

// ControllerCreator is a delegate to the creation method of a controller
// that is mapped as the primary route. (E.g. site.com/CONTROLLER is mapped
// to a NewXController method that implements this signature)
type ControllerCreator func(*http.Request) IController

// RouteManager provides the basic http request pipeline of the
// mvcapp framework
type RouteManager struct {
	SessionIDKey      string
	DefaultController string
	DefaultAction     string

	Routes         []*RouteMap
	SessionManager *SessionManager
}

// NewRouteManager returns a new route manager object with default
// controller and action tokens set to "Home" and "Index"
func NewRouteManager() *RouteManager {
	return &RouteManager{
		SessionIDKey: "MvcApp.SessionID",

		DefaultController: "Home",
		DefaultAction:     "Index",

		Routes:         make([]*RouteMap, 0),
		SessionManager: NewSessionManager(),
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

			// Set Browser Session
			browserSessionCookie, err := request.Cookie(manager.SessionIDKey)
			browserSessionID := ""

			if err != nil {
				browserSessionID = str.Random(32)
			} else {
				browserSessionID = browserSessionCookie.Value
			}

			browserSession := manager.SessionManager.GetSession(browserSessionID)
			controller.Session = browserSession

			// Prepare result
			result := icontroller.Execute(actionName, params)

			// Write controllers cookies
			for _, cookie := range controller.Cookies {
				http.SetCookie(response, cookie)
			}

			// Execute the response and return
			result.Execute(response)
			return
		}
	}
}
