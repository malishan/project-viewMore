package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"
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

	loglib.GenericFatalLog(apicontext.CustomContext{}, http.ListenAndServe(":"+port, handlers), nil)
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

//ErrorResponse - Return generic error message with error logging
func ErrorResponse(ctx apicontext.CustomContext, w http.ResponseWriter, responseErrorMessage string, statusCode int, logError error, fields loglib.FieldsMap) {
	var buf = new(bytes.Buffer)
	encoder := json.NewEncoder(buf)

	if logError != nil {
		loglib.GenericError(ctx, logError, fields)
	}

	encoder.Encode(InvalidParamResponse{Message: responseErrorMessage})
	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}

// HTTPResponse writes the HTTPResponse and renders the json: Uses context
func HTTPResponse(ctx apicontext.CustomContext, w http.ResponseWriter, statusCode int, msg string, data interface{}) {
	renderer := render.New()
	requestID := ctx.RequestID

	res := Response{}
	res.Meta.Code = statusCode
	res.Meta.Msg = msg
	res.Meta.RequestID = requestID
	res.Data = data
	renderer.JSON(w, statusCode, res)
}

//AddNoAuthRoute - Route without any Auth
func AddNoAuthRoute(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		MethodType:  methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, logRequest),
	}

	routes = append(routes, r)
}

//AddLoginRoute is to create route for login purpose only
func AddLoginRoute(methodName string, methodType string, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		MethodType:  methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, loginContextWrapper, logRequest),
	}

	routes = append(routes, r)
}

// AddRoute is to create routes with auth and user details
func AddRoute(methodName, methodType, mRoute string, handlerFunc http.HandlerFunc) {
	r := route{
		Name:        methodName,
		MethodType:  methodType,
		Pattern:     mRoute,
		HandlerFunc: useMiddleware(handlerFunc, validateContext, logRequest, createContext),
	}

	routes = append(routes, r)
}
