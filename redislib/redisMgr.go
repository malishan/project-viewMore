package redislib

import (
	"time"
)

// Get gets the value for given key
func Get(key string) (string, error) {
	client, connErr := getConnection()
	if connErr != nil {
		return "", connErr
	}

	cmd := client.Get(key)
	rslt, err := cmd.Result()

	return rslt, err
}

//Set sets the value for given key
func Set(key, value string) error {
	client, connErr := getConnection()
	if connErr != nil {
		return connErr
	}

	err := client.Set(key, value, 0).Err()

	return err
}

// SetWithExp is to set the key value with an expiry
func SetWithExp(expireSec int64, key string, value interface{}) error {
	client, connErr := getConnection()
	if connErr != nil {
		return connErr
	}

	err := client.SetNX(key, value, time.Duration(expireSec)*time.Second).Err()

	return err
}
