package mongolib

import (
	"context"
	"errors"
	"os"
	"project/project-viewMore/apicontext"
	"project/project-viewMore/loglib"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client *mongo.Client
	host   string
)

func init() {
	ctx := apicontext.CustomContext{}

	host = os.Getenv("MongoHost")
	if len(host) == 0 {
		loglib.GenericError(ctx, errors.New("MongoHost variable not set"), nil)
		return
	}

	loglib.GenericInfo(ctx, "mongo endpoint available: "+host, nil)
}

func connect() error {

	// creds := &options.Credential{
	// 	AuthSource: "admin",
	// 	Username:   "Alishan", //to-do: set username and password
	// 	Password:   "123456",
	// }

	servSelecTimeout := time.Duration(15) * time.Second
	connTimeout := time.Duration(10) * time.Second
	idleTime := time.Duration(2) * time.Minute
	socketTimeout := time.Duration(10) * time.Second
	appName := "View-More"
	maxPooling := uint64(0)

	clientOptions := &options.ClientOptions{
		AppName: &appName,
		//Auth:                   creds,
		ServerSelectionTimeout: &servSelecTimeout, //mongo driver will wait to select a server for an operation before giving up and raising an error
		ConnectTimeout:         &connTimeout,      //the driver will wait before a new connection attempt is aborted
		MaxConnIdleTime:        &idleTime,
		//Hosts:                  []string{"mongodb://localhost:27018"},
		MaxPoolSize:   &maxPooling,
		SocketTimeout: &socketTimeout, //the number of milliseconds a send or receive on a socket can take before timeout.
	}

	clientOptions = clientOptions.ApplyURI(host)

	err := clientOptions.Validate()
	if err != nil {
		return err
	}

	client, err = mongo.NewClient(clientOptions)
	if err != nil {
		return err
	}

	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	return nil
}

func getConnection() (*mongo.Client, error) {
	ctx := apicontext.CustomContext{}
	if client == nil {
		loglib.GenericInfo(ctx, "mongo connection broken, trying to reconnect...", nil)
		if err := connect(); err != nil {
			return nil, err
		}
	}

	return client, nil
}
