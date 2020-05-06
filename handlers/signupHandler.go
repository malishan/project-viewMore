package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/constant"
	"project/project-viewMore/core"
	"project/project-viewMore/mongolib"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

//UserSignUp : handler for user registration
func UserSignUp(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	multipartError := r.ParseMultipartForm(32 << 20)
	if multipartError != nil {
		core.ErrorResponse(ctx, w, "failed multipart parsing", http.StatusInternalServerError, fmt.Errorf("multipart parsing failed: %v", multipartError), nil)
		return
	}

	var user UserInfo

	formInputs := r.FormValue(constant.RegistrationForm)
	if len(formInputs) == 0 {
		core.ErrorResponse(ctx, w, "failed to parseForm", http.StatusInternalServerError, fmt.Errorf("form parsing failed: %v", multipartError), nil)
		return
	}

	unmarshalErr := json.Unmarshal([]byte(formInputs), &user)
	if unmarshalErr != nil {
		core.ErrorResponse(ctx, w, "failed to unmarshal input", http.StatusBadRequest, fmt.Errorf("failed to unmarshal input: %v", unmarshalErr), nil)
		return
	}

	if user.Name == "" || user.Email == "" || user.Password == "" {
		core.ErrorResponse(ctx, w, "user input missing", http.StatusBadRequest, errors.New("user input missing"), nil)
		return
	}

	var collection string
	if ctx.RoleID == constant.AdminRole {
		collection = constant.MongoAdminCollection
	} else if ctx.RoleID == constant.UserRole {
		collection = constant.MongoUserCollection
	} else {
		core.ErrorResponse(ctx, w, "invalid user role", http.StatusBadRequest, errors.New("invalid user role"), nil)
		return
	}

	filter := bson.M{"$or": []interface{}{
		bson.M{"name": user.Name},
		bson.M{"email": user.Email},
	}}

	isPresent, dbErr := mongolib.Exist(constant.MongoDatabaseName, collection, filter)
	if dbErr != nil {
		core.ErrorResponse(ctx, w, "failed to connect db", http.StatusInternalServerError, errors.New("failed to connect db"), nil)
		return
	} else if isPresent {
		core.ErrorResponse(ctx, w, "user already exists", http.StatusBadRequest, errors.New("user already exists"), nil)
		return
	}

	hashPwd, bcryptErr := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	if bcryptErr != nil {
		core.ErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError, fmt.Errorf("internal server error, %v", bcryptErr), nil)
		return
	}

	user.Password = string(hashPwd)

	insertRslt, insertErr := mongolib.CreateOne(constant.MongoDatabaseName, collection, user)
	if insertErr != nil {
		core.ErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError, fmt.Errorf("failed to create user, err: %v", insertErr), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "registration successfull", insertRslt.InsertedID)
}
