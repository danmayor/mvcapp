package mvcapp

import "net/http"

type RouteMap struct {
	ControllerName string
	Controller     interface{}
}

func NewRouteMap(name string, controller interface{}) RouteMap {
	return RouteMap{
		ControllerName: name,
		Controller:     &controller,
	}
}

func (routeMap *RouteMap) ExecuteAction(response http.ResponseWriter, request *http.Request) {
	// reflect over the controller in this map and execute the corresponding action method
}

func (routeMap *RouteMap) ExecuteResult(response http.ResponseWriter, actionResult interface{}) {
	// Called at the end of ExecuteAction to render the action result to the browser
}
