package main

import (
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"
	"project/project-viewMore/server"
)

const (
	version = "0.0.1"
)

func main() {
	m := loglib.FieldsMap{
		"Port":     server.Port,
		"SubRoute": server.SubRoute,
	}

	loglib.GenericInfo(apicontext.CustomContext{}, "STARTING VIEWMORE SERVER, VERSION: "+version, m)

	server.StartRoutes()
}
