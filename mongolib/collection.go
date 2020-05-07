package mongolib

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateOne - inserts data into mongo database
func CreateOne(database, collectionName string, d interface{}) (*mongo.InsertOneResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	insertOneRslt, err := collection.InsertOne(context.TODO(), d)
	if err != nil {
		return nil, err
	}

	return insertOneRslt, nil
}

// ReadOne - reads single document from mongo database
func ReadOne(database, collectionName string, filter, data interface{}) error {
	client, err := getConnection()
	if err != nil {
		return err
	}

	collection := client.Database(database).Collection(collectionName)
	err = collection.FindOne(context.TODO(), filter).Decode(data)
	if err != nil {
		return err
	}

	return nil
}

// Update - updates data into mongo database
func Update(database, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	updateResult, err := collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

//Exist verifies if document is present or not
//if it returns error then there is a connection error else boolean value specifies whether doc is present or not
func Exist(database, collectionName string, filter interface{}) (bool, error) {
	client, err := getConnection()
	if err != nil {
		return false, err
	}

	var i interface{}

	collection := client.Database(database).Collection(collectionName)
	err = collection.FindOne(context.TODO(), filter).Decode(&i)
	if err != nil {
		return false, nil
	}

	return true, nil
}

// AggregateAll executes aggregation query on a collection
//query []bson.M, data is a pointer to an array
func AggregateAll(database, collectionName string, query, data interface{}, options ...*options.AggregateOptions) error {
	client, err := getConnection()
	if err != nil {
		return err
	}

	collection := client.Database(database).Collection(collectionName)

	cursor, err := collection.Aggregate(context.TODO(), query, options...)
	if err != nil {
		return err
	}

	err = cursor.All(context.TODO(), data)
	return err
}
