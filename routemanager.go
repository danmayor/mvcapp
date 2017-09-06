package mvcapp

type RouteManager struct {
	Routes []RouteMap
}

func NewRouteManager() RouteManager {
	return RouteManager{
		Routes: make([]RouteMap, 0),
	}
}

func (manager *RouteManager) RegisterController(name string, controller IController) {
	manager.Routes = append(manager.Routes, NewRouteMap(name, controller))
}
