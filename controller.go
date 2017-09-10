/*
	Digivance MVC Application Framework
	Base Controller Feature
	Dan Mayor (dmayor@digivance.com)

	This file defines the base controller system functionality.
*/

package mvcapp

import (
	"net/http"

	"github.com/Digivance/str"
)

// IController defines the RegisterAction and Execute methods that
// need to be implemented by all controllers
type IController interface {
	RegisterAction(string, string, ActionMethod)
	Execute(string, []string) IActionResult
}

// Controller contains the basic members shared by custom controllers,
// also defines the required RegisterAction and Execute methods (below)
type Controller struct {
	IController

	Session      *Session
	Request      *http.Request
	ActionRoutes []*ActionMap
}

// RegisterAction allows package caller to map a controller action method to
// a given Http Request verb and action name (E.g. site.com/Controller/ActionName)
func (controller *Controller) RegisterAction(verb string, name string, method ActionMethod) {
	controller.ActionRoutes = append(controller.ActionRoutes, NewActionMap(verb, name, method))
}

// Execute the route manager to call the action method mapped to the provided
// actionName. params is the remainder of the url split by / represented as strings
func (controller *Controller) Execute(actionName string, params []string) IActionResult {
	for _, actionMethod := range controller.ActionRoutes {
		if str.Compare(actionMethod.Name, actionName) {
			return actionMethod.Method(params)
		}
	}

	return NewActionResult([]byte{})
}
