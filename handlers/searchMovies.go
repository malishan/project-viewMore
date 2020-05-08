package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/constant"
	"project/project-viewMore/core"
	"project/project-viewMore/mongolib"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

//SearchMovie handler fetches all movie info
func SearchMovie(w http.ResponseWriter, r *http.Request) {
	ctx := apicontext.UpgradeContext(r.Context())

	if ctx.RoleID != "" {
		core.ErrorResponse(ctx, w, "only guest users allowed", http.StatusForbidden, errors.New("only guest users allowed"), nil)
		return
	}

	movieName := r.URL.Query().Get(constant.MovieNameQueryString)

	if len(movieName) == 0 {
		core.ErrorResponse(ctx, w, "movieName missing from query param", http.StatusBadRequest, errors.New("movieName missing from query param"), nil)
		return
	}

	var movieResponse MovieDescription

	readErr := mongolib.ReadOne(constant.MongoDatabaseName, constant.MongoMovieCollection, bson.M{"name": movieName}, &movieResponse)
	if readErr != nil {
		if !strings.Contains(readErr.Error(), "no documents") {
			core.ErrorResponse(ctx, w, "failed to connect db", http.StatusInternalServerError, fmt.Errorf("failed to connect db, err: %v", readErr), nil)
		} else {
			core.ErrorResponse(ctx, w, "movieName not found", http.StatusBadRequest, fmt.Errorf("movie not found, err: %v", readErr), nil)
		}

		return
	}

	movieResponse.InsertTimestamp = 0
	movieResponse.RemoteAddress = ""
	movieResponse.UserAgent = ""
	movieResponse.UserEmail = ""

	core.HTTPResponse(ctx, w, http.StatusOK, "movie fetched successfully", movieResponse)
}
