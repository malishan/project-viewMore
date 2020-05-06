package core

import "net/http"

// Routes - list of route
type Routes []route

type route struct {
	Name                   string
	MethodType             string
	Pattern                string
	ResourcesPermissionMap map[string]uint8
	HandlerFunc            http.HandlerFunc
}
