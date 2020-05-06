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

//Response - complete structure of  HTTP response Meta + Data
type Response struct {
	Meta MetaData    `json:"meta"`
	Data interface{} `json:"data,omitempty"`
}

// MetaData of HTTP API response
type MetaData struct {
	Code      int    `json:"code"`
	Msg       string `json:"msg"`
	RequestID string `json:"requestId,omitempty"`
}
