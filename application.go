/*
	Digivance MVC Application Framework
	Global Application Scope Struct
	Dan Mayor (dmayor@digivance.com)
*/

package mvcapp

// Application is our global scope object (E.g. application wide configuration
// and manager systems)
type Application struct {
	// SessionKey is the Cookie Name used to store a connections
	// session ID
	SessionKey string

	// Sessions is our HTTP Session Manager system
	Sessions *SessionManager

	// Routes is our Route Manager system
	Routes *RouteManager
}

// NewApplication returns a new default MVC Application object
func NewApplication() Application {
	return Application{
		SessionKey: "SessionID",
		Routes:     NewRouteManager(),
	}
}
