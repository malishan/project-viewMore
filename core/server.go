package core

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
)

const (
	RoleID     = "roleID"
	UserID     = "userID"
	UserName   = "userName"
	UserEmail  = "email"
	TimeZone   = "timezone"
	Locale     = "locale"
	RequestID  = "requestId"
	timeFormat = time.RFC3339
)

var routes = make(Routes, 0)

//StartServer - http servers
func StartServer(port, subroute string) {

	allowedOrigins := handlers.AllowedOrigins([]string{"*"})

	allowedHeaders := handlers.AllowedHeaders([]string{
		"X-Requested-With",
		"X-CSRF-Token",
		"X-Auth-Token",
		"Content-Type",
		"processData",
		"contentType",
		"Origin",
		"Authorization",
		"Accept",
		"Client-Security-Token",
		"Accept-Encoding",
		TimeZone,
		Locale,
		RoleID,
		RequestID,
		UserID,
		UserName,
		UserEmail,
	})

	allowedMethods := handlers.AllowedMethods([]string{
		"POST",
		"GET",
		"DELETE",
		"PUT",
		"PATCH",
		"OPTIONS"})

	allowCredential := handlers.AllowCredentials()

	handlers := handlers.CORS(allowedHeaders, allowedMethods, allowedOrigins, allowCredential)(context.ClearHandler(newRouter(subroute)))

	log.Fatal(http.ListenAndServe(":"+port, handlers))
}

// NewRouter provides a mux Router.
// Handles all incoming request who matches registered routes against the request.
func newRouter(subroute string) *mux.Router {
	muxRouter := mux.NewRouter().StrictSlash(true)
	subRouter := muxRouter.PathPrefix(subroute).Subrouter()
	for _, route := range routes {
		subRouter.
			Methods(route.MethodType).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return muxRouter
}

// useMiddleware applies chains of middleware (ie: log, contextWrapper, validateAuth) handler into incoming request
// For example, logging middleware might write the incoming request details to a log
// Note - It applies in reverse order
func useMiddleware(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

// Ping API used to check the status of the service
func Ping(response http.ResponseWriter, request *http.Request) {
	HTTPPingResponse(response, http.StatusOK, map[string]string{"ping": "pong"})
}

// HTTPPingResponse writes the HTTPResponse and renders the json: Uses context
func HTTPPingResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	renderer := render.New()
	renderer.JSON(w, statusCode, data)
}

//AddNoAuthRoutes - Route without any Auth
func AddNoAuthRoutes(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		MethodType:  methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, logRequest),
	}

	routes = append(routes, r)
}
