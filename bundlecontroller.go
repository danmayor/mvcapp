package mvcapp

import "net/http"

type BundleController struct {
	*Controller
}

func NewBundleController(request *http.Request) *BundleController {
	return &BundleController{
		Controller: NewBaseController(request),
	}
}
