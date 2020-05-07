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
		"Add New Movie",
		http.MethodPost,
		"/add-movie",
		handlers.AddMovie,
	)

	/*---------------------------------Rate A Movie-------------------------------*/
	core.AddRoute(
		"Rate A Movie",
		http.MethodPost,
		"/rate-movie",
		handlers.AddMovieRating,
	)

	/*---------------------------------Comment on A Movie-------------------------------*/
	core.AddRoute(
		"Comment On A Movie",
		http.MethodPost,
		"/comment-movie",
		handlers.AddMovieComment,
	)

	/*---------------------------------Search Movie-------------------------------*/
	core.AddRoute(
		"Search Movie Info",
		http.MethodGet,
		"/search-movie",
		handlers.SearchMovie,
	)

	core.StartServer(Port, SubRoute)
}
