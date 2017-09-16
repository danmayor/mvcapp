/*
	Digivance MVC Application Framework
	Route Map Object
	Dan Mayor (dmayor@digivance.com)

	This file defines the generic MVC controller route maps
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
