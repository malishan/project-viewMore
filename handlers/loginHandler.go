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
	"project/project-viewMore/redislib"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// UserLogin : handler for user login
func UserLogin(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	err := r.ParseForm()
	if err != nil {
		core.ErrorResponse(ctx, w, "failed to parseForm", http.StatusInternalServerError, fmt.Errorf("file parsing failed: %v", err), nil)
		return
	}

	formInput := r.FormValue(constant.LoginForm)

	var userLogin UserInfo
	unmarshalErr := json.Unmarshal([]byte(formInput), &userLogin)
	if unmarshalErr != nil {
		core.ErrorResponse(ctx, w, "failed to unmarshal input", http.StatusInternalServerError, fmt.Errorf("input unmarshalling failed: %v", unmarshalErr), nil)
		return
	}

	var rsltUser struct {
		UserID   primitive.ObjectID `bson:"_id"`
		Password string             `bson:"pswd"`
	}

	var query bson.M
	if len(userLogin.Email) > 0 {
		query = bson.M{"email": userLogin.Email}
	} else {
		query = bson.M{"name": userLogin.Name}
	}

	if ctx.RoleID == constant.AdminRole {
		err = mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoAdminCollection, query, &rsltUser)
	} else if ctx.RoleID == constant.UserRole {
		err = mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoUserCollection, query, &rsltUser)
	} else {
		err = errors.New("user role not authentic")
	}

	if err != nil {
		core.ErrorResponse(ctx, w, "user not found", http.StatusBadRequest, fmt.Errorf("db operation failed: %v", err), nil)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(rsltUser.Password), []byte(userLogin.Password))
	if err != nil {
		core.ErrorResponse(ctx, w, "incorrect password", http.StatusBadRequest, fmt.Errorf("login failed: %v", err), nil)
		return
	}

	err = redislib.SetWithExp(constant.LoginExpirationTime, rsltUser.UserID.String(), time.Now().Unix())
	if err != nil {
		core.ErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError, fmt.Errorf("internal server error: %v", err), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "login successful", rsltUser.UserID.String())
}
