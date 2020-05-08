package core

import (
	"errors"
	"fmt"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/constant"
	"project/project-viewMore/loglib"
	"project/project-viewMore/mongolib"

	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
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

// createContext generates context for incoming request and appends in request Context
func createContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header
		ctx := r.Context()

		requestID := header.Get(RequestID)
		if requestID == "" {
			requestID = uuid.NewV4().String()
		}

		userEmail, userName, roleID := header.Get(UserEmail), header.Get(UserName), header.Get(RoleID)

		apiCtx := apicontext.APIContext{
			RequestID: requestID,
			UserName:  userName,
			Email:     userEmail,
			RoleID:    roleID,
		}

		ctx = apicontext.WithAPIContext(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateContext incoming request context
func validateContext(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		apiCtx, ok := ctx.Value(apicontext.APICtx).(apicontext.APIContext)

		tempContext := apicontext.CustomContext{}
		tempContext.APIContext = apiCtx

		if !ok {
			ErrorResponse(tempContext, w, "internal server error, cannot cast context", http.StatusInternalServerError, errors.New("context parsing failed"), nil)
			return
		}

		if r.URL.Path == "/viewmore/search-movie" {
			next.ServeHTTP(w, r)
			return
		} else if !(apiCtx.RoleID == constant.AdminRole || apiCtx.RoleID == constant.UserRole) {
			ErrorResponse(tempContext, w, "roleID header incorrect", http.StatusBadRequest, errors.New("roleID header incorrect"), nil)
			return
		}

		if len(apiCtx.Email) == 0 {
			ErrorResponse(tempContext, w, "email header missing", http.StatusBadRequest, errors.New("email header missing"), nil)
			return
		}

		var userID struct {
			ID primitive.ObjectID `json:"_id" bson:"_id"`
		}

		collectionName := ""
		if apiCtx.RoleID == constant.AdminRole {
			collectionName = constant.MongoAdminCollection
		} else {
			collectionName = constant.MongoUserCollection
		}

		err := mongolib.ReadOne(constant.MongoDatabaseName, collectionName, bson.M{"email": apiCtx.Email}, &userID)
		if err != nil {
			ErrorResponse(tempContext, w, "email incorrect, please register (or mongo inactive)", http.StatusBadRequest, fmt.Errorf("failed to fetch userID, err: %v", err), nil)
			return
		}
		apiCtx.UserID = userID.ID.Hex()

		if len(apiCtx.UserID) == 0 {
			ErrorResponse(tempContext, w, "user not registered", http.StatusBadRequest, errors.New("user not registered"), nil)
			return
		}

		ctx = apicontext.WithAPIContext(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
