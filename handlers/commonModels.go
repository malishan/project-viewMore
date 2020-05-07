package handlers

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
	ID              string        `json:"_id,omitempty" bson:"_id,omitempty"`
	Name            string        `json:"name,omitempty" bson:"name,omitempty"`
	AvgRating       int           `json:"avgRating,omitempty" bson:"avgRating,omitempty"`
	RatingCount     int           `json:"ratingCount,omitempty" bson:"ratingCount,omitempty"`
	InsertTimestamp int64         `json:"timestamp,omitempty" bson:"timestamp,omitempty"`
	UserAgent       string        `json:"userAgent,omitempty" bson:"userAgent,omitempty"`
	RemoteAddress   string        `json:"remoteAddr,omitempty" bson:"remoteAddr,omitempty"`
	Comments        []CommentInfo `json:"comments,omitempty" bson:"comments,omitempty"`
}

//CommentInfo - comment info by each user on a particular movie
type CommentInfo struct {
	UserName    string   `json:"userName,omitempty" bson:"userName,omitempty"`
	UserEmail   string   `json:"userEmail,omitempty" bson:"userEmail,omitempty"`
	CommentTime string   `json:"commentTime,omitempty" bson:"commentTime,omitempty"`
	UserRating  int      `json:"userRating,omitempty" bson:"userRating,omitempty"`
	UserComment []string `json:"userComment,omitempty" bson:"userComment,omitempty"`
}
