/*
	Digivance MVC Application Framework
	Route Mapping Features
	Dan Mayor (dmayor@digivance.com)

	This file defines a generic "Route Map". Route maps are intended to define the controller
	object to be used for requests to the provided controller section of the requested url.
*/

package mvcapp

// RouteMap is used to map the controller portion of the requested URL
// to a controller struct that implements IController
type RouteMap struct {
	// ControllerName is controller portion of the url that this route map responds to
	ControllerName   string

	// CreateController is the New*Controller method we call to invoke an instance of
	// the core controller object (E.g. custom controllers simply provide and register
	// a method to this map)
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
