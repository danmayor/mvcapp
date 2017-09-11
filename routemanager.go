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

// RouteManager provides the basic http request pipeline of the
// mvcapp framework
type RouteManager struct {
	DefaultController string
	DefaultAction     string

	Routes []*RouteMap
}

// NewRouteManager returns a new route manager object with default
// controller and action tokens set to "Home" and "Index"
func NewRouteManager() *RouteManager {
	return &RouteManager{
		DefaultController: "Home",
		DefaultAction:     "Index",
		Routes:            make([]*RouteMap, 0),
	}
}

// RegisterController is used to map a custom controller object to the
// controller section of the requested url (E.g. "site.com/CONTROLLER/action")
func (manager *RouteManager) RegisterController(name string, controller IController) {
	manager.Routes = append(manager.Routes, NewRouteMap(name, controller))
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
			route.Controller.SetRequest(request)
			// Set Cookies
			// Set Session

			result := route.Controller.Execute(actionName, params)
			result.Execute(response)
			return
		}
	}
}
