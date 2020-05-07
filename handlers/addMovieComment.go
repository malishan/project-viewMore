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

//AddMovieComment handler adds comment to a movie by loggedIn users only
func AddMovieComment(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	if ctx.RoleID != constant.UserRole {
		core.ErrorResponse(ctx, w, "user not allowed", http.StatusForbidden, errors.New("user not allowed"), nil)
		return
	} else if ctx.UserID == "" {
		core.ErrorResponse(ctx, w, "userID missing from header", http.StatusBadRequest, errors.New("userID missing from header"), nil)
		return
	}

	_, err := redislib.Get(ctx.UserID)
	if err != nil {
		core.ErrorResponse(ctx, w, "user login required", http.StatusBadRequest, fmt.Errorf("user not loggedIn, err: %v", err), nil)
		return
	}

	var usrComm UserRatingAndComment
	decodeErr := json.NewDecoder(r.Body).Decode(&usrComm)
	if decodeErr != nil {
		core.ErrorResponse(ctx, w, "failed to decode input", http.StatusBadRequest, fmt.Errorf("input decoding failed: %v", decodeErr), nil)
		return
	}

	if usrComm.MovieName == "" || usrComm.Comment[0] == "" {
		core.ErrorResponse(ctx, w, "incomplete request body", http.StatusBadRequest, errors.New("movieName or comment missing from request body"), nil)
		return
	}

	isPresent, existErr := mongolib.Exist(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": usrComm.MovieName})
	if existErr != nil {
		core.ErrorResponse(ctx, w, "failed to connect mongo", http.StatusInternalServerError, fmt.Errorf("failed to connect mongo: %v", existErr), nil)
		return
	} else if !isPresent {
		core.ErrorResponse(ctx, w, "movieName not found", http.StatusBadRequest, errors.New("movieName not found"), nil)
		return
	}

	isPresent, existErr = mongolib.Exist(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": usrComm.MovieName, "userFeedback.userID": ctx.UserID})
	if existErr != nil {
		core.ErrorResponse(ctx, w, "failed to connect mongo", http.StatusBadRequest, fmt.Errorf("failed to connect mongo: %v", existErr), nil)
		return
	}

	var update, filter bson.M
	if isPresent {
		filter = bson.M{"name": usrComm.MovieName, "userFeedback.userID": ctx.UserID}
		update = bson.M{
			"$set": bson.M{
				"userFeedback.$.updatedTime": time.Now().UnixNano() / int64(time.Millisecond),
				"$addToSet":                  bson.M{"userFeedback.userComment": usrComm.Comment},
			}}
	} else {

		feedback := UserFeedback{
			UserID:      ctx.UserID,
			UserEmail:   ctx.Email,
			UpdatedTime: time.Now().UnixNano() / int64(time.Millisecond),
			UserComment: []string{usrComm.Comment[0]},
		}

		filter = bson.M{"name": usrComm.MovieName}
		update = bson.M{
			"$set": bson.M{
				"$addToSet": bson.M{
					"userFeedback": feedback},
			}}
	}

	updateRslt, updateErr := mongolib.Update(constant.MongoDatabaseName, constant.MongoMovieCollection, filter, update)

	if updateErr != nil {
		core.ErrorResponse(ctx, w, "comment update failed", http.StatusBadRequest, fmt.Errorf("comment update failed: %v", updateErr), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "movie comment added successfully", fmt.Sprintf("total updated document: %v", updateRslt.ModifiedCount))
}
