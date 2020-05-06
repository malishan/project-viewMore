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

// CreateMany - inserts many data into mongo database
func CreateMany(database, collectionName string, d ...interface{}) (*mongo.InsertManyResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	insertManyRslt, err := collection.InsertMany(context.TODO(), d)
	if err != nil {
		return nil, err
	}

	return insertManyRslt, nil
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

// ReadAll - reads multiple documents from mongo database
func ReadAll(database, collectionName string, filter, data interface{}, opts ...*options.FindOptions) error {
	client, err := getConnection()
	if err != nil {
		return err
	}
	var findOptions *options.FindOptions
	if len(opts) > 0 {
		findOptions = opts[0]
	}

	collection := client.Database(database).Collection(collectionName)
	cursor, err := collection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return err
	}
	defer cursor.Close(context.TODO())

	err = cursor.All(context.TODO(), data)
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

// UpdateAll - updates multiple documents into mongo database
func UpdateAll(database, collectionName string, filter, update interface{}) (*mongo.UpdateResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	updateResult, err := collection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}

	return updateResult, nil
}

// Delete - removes single doc data from the database
func Delete(database, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

// DeleteAll - removes all doc data from the database
func DeleteAll(database, collectionName string, filter interface{}) (*mongo.DeleteResult, error) {
	client, err := getConnection()
	if err != nil {
		return nil, err
	}

	collection := client.Database(database).Collection(collectionName)
	deleteResult, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	return deleteResult, nil
}

//CountDocuments returns document count of a collection
func CountDocuments(database, collectionName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	client, err := getConnection()
	if err != nil {
		return 0, err
	}

	var countOptions *options.CountOptions
	if len(opts) > 0 {
		countOptions = opts[0]
	}

	collection := client.Database(database).Collection(collectionName)
	count, err := collection.CountDocuments(context.TODO(), filter, countOptions)
	if err != nil {
		return 0, err
	}

	return count, nil
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

// GetDistinct gets the distinct values for the field name provided
func GetDistinct(database, collectionName, fieldName string, filter interface{}) (interface{}, error) {
	client, err := getConnection()
	if err != nil {
		return false, err
	}
	collection := client.Database(database).Collection(collectionName)
	result, err := collection.Distinct(context.TODO(), fieldName, filter, nil)
	if err != nil {
		return nil, err
	}
	return result, nil
}
