/*
	Digivance MVC Application Framework
	Action Map Features
	Dan Mayor (dmayor@digivance.com)

	This file defines functionality for mapping an action method to an http request optionally
	boud to an http verb.
*/

package mvcapp

// ActionMethod defines the method signature for controller action methods
type ActionMethod func([]string) *ActionResult

// ActionMap is used to define the HTTP Verb, Controller's Action Name
// and the corresponding action method
type ActionMap struct {
	// Verb is the HTTP Verb to bind this mapping to, blank to respond to all
	Verb string

	// Name is the site.com/controller/<ACTION> name that this map handles
	Name string

	// Method is the actual action method to execute on the controller
	Method ActionMethod
}

// NewActionMap returns a new ActionMap struct populated with the given parameters
func NewActionMap(httpVerb string, actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   httpVerb,
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewGetActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to get
func NewGetActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "GET",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewPostActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to post
func NewPostActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "POST",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewPutActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to put
func NewPutActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "PUT",
		Name:   actionName,
		Method: actionMethod,
	}
}

// NewDeleteActionMap returns a new ActionMap struct populated with the given parameters
// and sets the HTTP Verb to delete
func NewDeleteActionMap(actionName string, actionMethod ActionMethod) *ActionMap {
	return &ActionMap{
		Verb:   "DELETE",
		Name:   actionName,
		Method: actionMethod,
	}
}
