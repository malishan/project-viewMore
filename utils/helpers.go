package utils

import (
	"project/project-viewMore/redislib"
)

// IsUserLoggedIn checks if user is logged In
func IsUserLoggedIn(userID string) error {
	_, getErr := redislib.Get(userID)
	return getErr
}
