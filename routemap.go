/*
	Digivance MVC Application Framework
	Route Mapping Features
	Dan Mayor (dmayor@digivance.com)

	This file defines a generic "Route Map". Route maps are intended to define the controller
	object to be used for requests to the provided controller section of the requested url.

	This package is released under as open source under the LGPL-3.0 which can be found:
	https://opensource.org/licenses/LGPL-3.0
*/

package mvcapp

// RouteMap is used to map the controller portion of the requested URL
// to a controller struct that implements IController
type RouteMap struct {
	ControllerName   string
	CreateController ControllerCreator
}

// NewRouteMap returns a new RouteMap object populated with the provided
// name and controller. (E.g. site.com/CONTROLLER/* is mapped to the
// provided controller creator)
func NewRouteMap(name string, creator ControllerCreator) *RouteMap {
	return &RouteMap{
		ControllerName:   name,
		CreateController: creator,
	}
}
