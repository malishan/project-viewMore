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
)

//AddMovie handler adds new movie info by the admin
func AddMovie(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	if ctx.RoleID != constant.AdminRole {
		core.ErrorResponse(ctx, w, "user not allowed", http.StatusForbidden, errors.New("user not allowed"), nil)
		return
	} else if ctx.UserID == "" {
		core.ErrorResponse(ctx, w, "userID missing from header", http.StatusBadRequest, errors.New("userID missing from header"), nil)
		return
	}

	_, err := redislib.Get(ctx.UserID)
	if err != nil {
		core.ErrorResponse(ctx, w, "admin login required", http.StatusBadRequest, fmt.Errorf("user not loggedIn, err: %v", err), nil)
		return
	}

	var movie MovieDescription
	decodeErr := json.NewDecoder(r.Body).Decode(&movie)
	if decodeErr != nil {
		core.ErrorResponse(ctx, w, "failed to decode input", http.StatusBadRequest, fmt.Errorf("input decoding failed: %v", decodeErr), nil)
		return
	}

	if movie.Name == "" {
		core.ErrorResponse(ctx, w, "movie name missing from request body", http.StatusBadRequest, errors.New("movie name missing from request body"), nil)
		return
	}

	isExist, dbErr := mongolib.Exist(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": movie.Name})
	if dbErr != nil {
		core.ErrorResponse(ctx, w, "failed to connect db", http.StatusInternalServerError, errors.New("failed to connect db"), nil)
		return
	} else if isExist {
		core.ErrorResponse(ctx, w, "movieName already exists", http.StatusBadRequest, errors.New("movieName already exists"), nil)
		return
	}

	movie.InsertTimestamp = time.Now().UnixNano() / int64(time.Millisecond)
	movie.RemoteAddress = r.RemoteAddr
	movie.UserAgent = r.Header.Get(constant.UserAgent)

	insertRslt, insertErr := mongolib.CreateOne(constant.MongoDatabaseName, constant.MongoMovieCollection, movie)
	if insertErr != nil {
		core.ErrorResponse(ctx, w, "internal server error", http.StatusInternalServerError, fmt.Errorf("failed to insert movie, err: %v", insertErr), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "movie added successfully", insertRslt.InsertedID)
}