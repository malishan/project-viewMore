package redislib

import (
	"errors"
	"os"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"

	"github.com/go-redis/redis"
)

var (
	client *redis.Client

	address string
)

func init() {
	ctx := apicontext.CustomContext{}

	address = os.Getenv("Redis_Host")
	if len(address) == 0 {
		loglib.GenericFatalLog(ctx, errors.New("env variable not set for redis address"), nil)
	}

	loglib.GenericInfo(ctx, "redis endpoint available: "+address, nil)
}

func connectRedis() error {

	clientOptions := &redis.Options{
		Addr: "localhost:6379",
	}

	client = redis.NewClient(clientOptions)

	_, err := client.Ping().Result()

	return err
}

func getConnection() (*redis.Client, error) {
	if client == nil {
		err := connectRedis()
		if err == nil {
			return client, nil
		}

		return nil, err
	}

	return client, nil
}
