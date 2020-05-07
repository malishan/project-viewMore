package handlers

import "sync"

// UserInfo forms the basic structure of an user
type UserInfo struct {
	UserID   string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"pswd" bson:"pswd"`
	Phone    []string `json:"phoneNo,omitempty" bson:"phoneNo,omitempty"`
	Address  string   `json:"address,omitempty" bson:"address,omitempty"`
}

//MovieDescription - info of each movie
type MovieDescription struct {
	ID              string         `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string         `json:"name,omitempty" bson:"name,omitempty"`
	AvgRating       float64        `json:"avgRating,omitempty" bson:"avgRating,omitempty"`
	RatingCount     int            `json:"ratingCount,omitempty" bson:"ratingCount,omitempty"`
	InsertTimestamp int64          `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	UserAgent       string         `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	RemoteAddress   string         `json:"remoteAddr,omitempty" bson:"remoteAddr,omitempty"`
	IndivFeedback   []UserFeedback `json:"userFeedback,omitempty" bson:"userFeedback,omitempty"`
}

//UserFeedback - comments and rating info by each user on a particular movie
type UserFeedback struct {
	UserID      string   `json:"userID,omitempty" bson:"userID,omitempty"`
	UserEmail   string   `json:"userEmail,omitempty" bson:"userEmail,omitempty"`
	UpdatedTime int64    `json:"updatedTime,omitempty" bson:"updatedTime,omitempty"`
	UserRating  float64  `json:"userRating,omitempty" bson:"userRating,omitempty"`
	UserComment []string `json:"userComment,omitempty" bson:"userComment,omitempty"`
}

//UserRatingAndComment - rating and comment of individual user
type UserRatingAndComment struct {
	MovieName string  `json:"movieName" bson:"movieName"`
	Rating    float64 `json:"rating" bson:"rating"`
	Comment   string  `json:"comment" bson:"comment"`
}

var (
	//Mutex - global mutext to handle synchronization
	Mutex sync.Mutex
)
