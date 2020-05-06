package main

import (
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"
)

const (
	version = "0.0.1"
)

func main() {
	loglib.GenericInfo(apicontext.CustomContext{}, "STARTING VIEWMORE SERVER, VERSION: "+version, nil)
}
