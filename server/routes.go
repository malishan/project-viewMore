package server

import (
	"net/http"
	"project/project-viewMore/core"
	"project/project-viewMore/handlers"
)

const (
	//Port is the port number of the current service
	Port = "8989"
	//SubRoute is the basic route of all APIs of this service
	SubRoute = "/viewmore"
)

func init() {

	//Adding ping api check for all the services for health check
	core.AddNoAuthRoutes(
		"Ping Check",
		http.MethodGet,
		"/health",
		core.Ping)
}

//StartRoutes is the starting point of the http service
func StartRoutes() {

	/*---------------------------------Registration-------------------------------*/
	core.AddLoginRoutes(
		"SignUp",
		http.MethodPost,
		"/register",
		nil,
		handlers.UserSignUp,
	)

	/*---------------------------------Login-------------------------------*/
	core.AddLoginRoutes(
		"Login",
		http.MethodPost,
		"/login",
		nil,
		handlers.UserLogin,
	)

	core.StartServer(Port, SubRoute)
}
