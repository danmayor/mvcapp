package mvcapp

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/Digivance/str"
)

type RouteManager struct {
	DefaultController string
	DefaultAction     string

	Routes []RouteMap
}

func NewRouteManager() *RouteManager {
	return &RouteManager{
		DefaultController: "Home",
		DefaultAction:     "Index",
		Routes:            make([]RouteMap, 0),
	}
}

func (manager *RouteManager) RegisterController(name string, controller IController) {
	manager.Routes = append(manager.Routes, NewRouteMap(name, controller))
}

func (manager *RouteManager) HandleRequest(response http.ResponseWriter, request *http.Request) {
	parts := str.Split(request.URL.EscapedPath(), '/')

	var controllerName string
	var actionName string

	if len(parts) >= 2 {
		controllerName = parts[0]
		actionName = parts[1]
	} else {
		if len(parts) == 1 {
			controllerName = parts[0]
			actionName = manager.DefaultAction
		} else {
			controllerName = manager.DefaultController
			actionName = manager.DefaultAction
		}
	}

	for _, v := range manager.Routes {
		if str.Compare(controllerName, v.ControllerName) {
			cv := v.Controller
			controller := reflect.ValueOf(cv).Elem()
			if controller != reflect.Zero(controller.Type()) {
				/*
					controllerAction := controller.FieldByNameFunc(func(name string) bool {
						return str.Compare(name, actionName)
					})
				*/

				controllerAction := runtime.FuncForPC(reflect.ValueOf(cv).Pointer()).Name
				fmt.Println(actionName)
				fmt.Println(controllerAction)

				/*
					if controllerAction != reflect.Zero(controllerAction.Type()) {
						controllerAction.Call([]reflect.Value{})
					}
				*/
			}
		}
	}
}
