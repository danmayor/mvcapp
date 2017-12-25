/*
	Digivance MVC Application Framework - Unit Tests
	Route Map Feature Tests
	Dan Mayor (dmayor@digivance.com)

	This file defines the version 0.2.0 compatibility of routemap.go functions. These functions are written
	to demonstrate and test the intended use cases of the functions in routemap.go
*/

package mvcapp_test

import (
	"net/http"
	"testing"

	"github.com/digivance/mvcapp"
)

// routeMapControllerCreator is used internally to test creating the new route map
func routeMapControllerCreator(request *http.Request) mvcapp.IController {
	return nil
}

// TestNewRouteMap ensures that the NewRouteMap method returns the expected value
func TestNewRouteMap(t *testing.T) {
	routeMap := mvcapp.NewRouteMap("Test", routeMapControllerCreator)
	if routeMap.ControllerName != "Test" || routeMap.CreateController == nil {
		t.Fatal("Nope")
	}
}
