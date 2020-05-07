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
	core.AddNoAuthRoute(
		"Ping Check",
		http.MethodGet,
		"/health",
		core.Ping)
}

//StartRoutes is the starting point of the http service
func StartRoutes() {

	/*---------------------------------Registration-------------------------------*/
	core.AddLoginRoute(
		"SignUp",
		http.MethodPost,
		"/register",
		handlers.UserSignUp,
	)

	/*---------------------------------Login-------------------------------*/
	core.AddLoginRoute(
		"Login",
		http.MethodPost,
		"/login",
		handlers.UserLogin,
	)

	/*---------------------------------Add New Movie-------------------------------*/
	core.AddRoute(
		"",
		http.MethodPost,
		"/add-movie",
		handlers.AddMovie,
	)

	core.StartServer(Port, SubRoute)
}
