package core

import (
	"errors"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"

	"github.com/twinj/uuid"
)

//logRequest logs each HTTP incoming Requests
func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		loglib.GenericTrace(apicontext.CustomContext{}, "API Called: "+r.URL.Path, nil)
		next.ServeHTTP(w, r)
	})
}

// loginContextWrapper incoming request context
func loginContextWrapper(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		tempContext := apicontext.CustomContext{}

		roleID := r.Header.Get(RoleID)
		if len(roleID) == 0 {
			loglib.GenericError(tempContext, errors.New("missing roleID header"), nil)
			ErrorResponse(tempContext, w, "userRole not defined", http.StatusBadRequest, errors.New("missing roleID header"), nil)
			return
		}

		requestID := r.Header.Get(RequestID)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}

		apiCtx := apicontext.APIContext{
			RoleID:    roleID,
			RequestID: requestID,
		}

		ctx = apicontext.WithAPIContext(ctx, apiCtx)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
