package core

import (
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"
)

//logRequest logs each HTTP incoming Requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loglib.GenericTrace(apicontext.CustomContext{}, "API Called: "+r.URL.Path, nil)
		next.ServeHTTP(w, r)
	})
}
