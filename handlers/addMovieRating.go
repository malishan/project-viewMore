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

//AddMovieRating handler add rating of a movie by loggedIn users only
func AddMovieRating(w http.ResponseWriter, r *http.Request) {
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

	var usrRating UserRatingAndComment
	decodeErr := json.NewDecoder(r.Body).Decode(&usrRating)
	if decodeErr != nil {
		core.ErrorResponse(ctx, w, "failed to decode input", http.StatusBadRequest, fmt.Errorf("input decoding failed: %v", decodeErr), nil)
		return
	}

	if usrRating.MovieName == "" || usrRating.Rating == 0 {
		core.ErrorResponse(ctx, w, "incomplete request body", http.StatusBadRequest, errors.New("movieName or rating missing from request body"), nil)
		return
	}

	Mutex.Lock()

	var prevRating struct {
		Rate  float64 `json:"avgRating,omitempty" bson:"avgRating,omitempty"`
		Count int     `json:"ratingCount,omitempty" bson:"ratingCount,omitempty"`
	}
	readErr := mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": usrRating.MovieName}, &prevRating)
	if readErr != nil {
		core.ErrorResponse(ctx, w, "movieName not found", http.StatusBadRequest, fmt.Errorf("failed to fetch movieDetails: %v", readErr), nil)
		Mutex.Unlock()
		return
	}

	newAvg := (prevRating.Rate*float64(prevRating.Count) + usrRating.Rating) / float64(prevRating.Count+1)

	isPresent, existErr := mongolib.Exist(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": usrRating.MovieName, "userFeedback.userID": ctx.UserID})
	if existErr != nil {
		core.ErrorResponse(ctx, w, "failed to connect mongo", http.StatusBadRequest, fmt.Errorf("failed to connect mongo: %v", existErr), nil)
		return
	}

	var update, filter bson.M
	if isPresent {
		filter = bson.M{"name": usrRating.MovieName, "userFeedback.userID": ctx.UserID}
		update = bson.M{
			"$set": bson.M{
				"avgRating":                  newAvg,
				"ratingCount":                prevRating.Count + 1,
				"userFeedback.$.userRating":  usrRating.Rating,
				"userFeedback.$.updatedTime": time.Now().UnixNano() / int64(time.Millisecond)}}
	} else {

		feedback := UserFeedback{
			UserID:      ctx.UserID,
			UserEmail:   ctx.Email,
			UserRating:  usrRating.Rating,
			UpdatedTime: time.Now().UnixNano() / int64(time.Millisecond),
		}

		filter = bson.M{"name": usrRating.MovieName}
		update = bson.M{
			"$set": bson.M{
				"avgRating":   newAvg,
				"ratingCount": prevRating.Count + 1,
				"$addToSet": bson.M{
					"userFeedback": feedback},
			}}
	}

	updateRslt, updateErr := mongolib.Update(constant.MongoDatabaseName, constant.MongoMovieCollection, filter, update)

	Mutex.Unlock()

	if updateErr != nil {
		core.ErrorResponse(ctx, w, "rating update failed", http.StatusBadRequest, fmt.Errorf("rating update failed: %v", updateErr), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "movie added successfully", fmt.Sprintf("total updated document: %v", updateRslt.ModifiedCount))
}
