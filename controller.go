package mvcapp

import "net/http"

// IController is used internally to make the custom controller objects
// visible to our RouteManager system
type IController interface {
}

// Controller contains the basic members shared by custom controllers
type Controller struct {
	IController

	Session *Session
	Request *http.Request
}
