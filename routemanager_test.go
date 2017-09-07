package mvcapp

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

type TestController struct {
	Controller
}

func (controller TestController) Index() {
	fmt.Println("Yay it works!")
}

func TestRouteManager(t *testing.T) {
	controller := TestController{}
	mgr := NewRouteManager()
	mgr.RegisterController("Test", controller)

	request := httptest.NewRequest("GET", "/Test/Index", nil)
	response := httptest.NewRecorder()

	mgr.HandleRequest(response, request)
}
