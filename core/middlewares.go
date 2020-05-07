package core

import (
	"errors"
	"fmt"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/constant"
	"project/project-viewMore/loglib"
	"project/project-viewMore/mongolib"
	"project/project-viewMore/redislib"
	"time"

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

		userEmail, userID, userName, roleID := header.Get(UserEmail), header.Get(UserID), header.Get(UserName), header.Get(RoleID)
		//token := header.Get(AppAPIToken)

		apiCtx := apicontext.APIContext{
			RequestID: requestID,
			UserID:    userID,
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

		if len(apiCtx.Email) == 0 {
			ErrorResponse(tempContext, w, "email header missing", http.StatusBadRequest, errors.New("email header missing"), nil)
			return
		}

		if !(len(apiCtx.RoleID) == 0 || apiCtx.RoleID == constant.AdminRole || apiCtx.RoleID == constant.UserRole) {
			ErrorResponse(tempContext, w, "email header missing", http.StatusBadRequest, errors.New("email header missing"), nil)
			return
		}

		if apiCtx.RoleID == constant.AdminRole || apiCtx.RoleID == constant.UserRole {
			go updateLoginStatus(tempContext)
		}

		ctx = apicontext.WithAPIContext(ctx, apiCtx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func updateLoginStatus(ctx apicontext.CustomContext) {

	var (
		err    error
		userID struct {
			ID primitive.ObjectID `json:"_id" bson:"_id"`
		}
	)

	if ctx.RoleID == constant.AdminRole {
		err = mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoAdminCollection, bson.M{"email": ctx.Email}, &userID)
	} else {
		err = mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoUserCollection, bson.M{"email": ctx.Email}, &userID)
	}

	if err != nil {
		loglib.GenericError(ctx, fmt.Errorf("mongo call failed, err: %d", err), nil)
	}

	err = redislib.SetWithExp(constant.LoginExpirationTime, userID.ID.String(), time.Now().Unix())
	if err != nil {
		loglib.GenericError(ctx, fmt.Errorf("redis call failed, err: %d", err), nil)
	}
}
