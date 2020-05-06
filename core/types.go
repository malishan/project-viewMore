package core

import "net/http"

// Routes - list of route
type Routes []route

type route struct {
	Name        string
	MethodType  string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//InvalidParamResponse returns the error response
type InvalidParamResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}
