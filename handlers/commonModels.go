package handlers

// UserInfo forms the basic structure of an user
type UserInfo struct {
	UserID   string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string   `json:"name" bson:"name"`
	Email    string   `json:"email" bson:"email"`
	Password string   `json:"pswd" bson:"pswd"`
	Phone    []string `json:"phoneNo" bson:"phoneNo"`
	Address  string   `json:"address" bson:"address"`
}
