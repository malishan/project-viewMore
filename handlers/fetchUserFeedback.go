package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/constant"
	"project/project-viewMore/core"
	"project/project-viewMore/mongolib"

	"go.mongodb.org/mongo-driver/bson"
)

//FetchUserFeedback handler returns all comments and ratings of a single loggedIn user
func FetchUserFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	if ctx.RoleID != constant.UserRole {
		core.ErrorResponse(ctx, w, "user not allowed", http.StatusForbidden, errors.New("user not allowed"), nil)
		return
	}

	// _, err := redislib.Get(ctx.UserID)
	// if err != nil {
	// 	core.ErrorResponse(ctx, w, "user login required", http.StatusBadRequest, fmt.Errorf("user not loggedIn, err: %v", err), nil)
	// 	return
	// }

	q1 := bson.M{"$unwind": "$userFeedback"}
	q2 := bson.M{"$group": bson.M{"_id": bson.M{"movieName": "$name", "userID": "$userFeedback.userID"}, "movieName": bson.M{"$first": "$name"}, "userID": bson.M{"$first": "$userFeedback.userID"}, "rating": bson.M{"$first": "$userFeedback.userRating"}, "comments": bson.M{"$first": "$userFeedback.userComment"}}}
	q3 := bson.M{"$match": bson.M{"userID": ctx.UserID}}
	q4 := bson.M{"$project": bson.M{"_id": 0, "movieName": 1, "rating": 1, "comments": 1}}

	query := []bson.M{q1, q2, q3, q4}

	var feedbackRsp []interface{}
	aggErr := mongolib.AggregateAll(constant.MongoDatabaseName, constant.MongoMovieCollection, query, &feedbackRsp)
	if aggErr != nil {
		core.ErrorResponse(ctx, w, "failed to fetch user feedback info", http.StatusBadRequest, fmt.Errorf("aggregate mongo query failed: %v", aggErr), nil)
		return
	}

	core.HTTPResponse(ctx, w, http.StatusOK, "user feedback fetched successfully", feedbackRsp)
}
