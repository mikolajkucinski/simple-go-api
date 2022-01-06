package internal

import (
	proto_files "awesomeProject/internal/proto-files"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DbConnector struct {
	client                             *mongo.Client
	UserCollection, EmployeeCollection *mongo.Collection
	Context                            context.Context
}

func (dbConnector *DbConnector) Connect() {
	var connectionError error

	dbConnector.Context, _ = context.WithTimeout(context.Background(), 10*time.Second)
	dbConnector.client, connectionError = mongo.Connect(dbConnector.Context, options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if connectionError != nil {
		panic(connectionError)
	}

	dbConnector.UserCollection = dbConnector.client.Database("go_test").Collection("UserCollection")
	dbConnector.EmployeeCollection = dbConnector.client.Database("go_test").Collection("EmployeeCollection")
}

func (dbConnector *DbConnector) FindUserByUserId(userId string) ([]bson.M, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return []bson.M{}, err
	}

	filterCursor, err := dbConnector.UserCollection.Find(dbConnector.Context, bson.M{"id": id})
	if err != nil {
		return []bson.M{}, err
	}

	var user []bson.M
	if err := filterCursor.All(dbConnector.Context, &user); err != nil {
		return []bson.M{}, err
	}

	return user, nil
}

func (dbConnector *DbConnector) InsertUser(user *proto_files.User) (string, error) {
	result, err := dbConnector.UserCollection.InsertOne(dbConnector.Context, bson.D{
		{Key: "firstName", Value: user.GetFirstName()},
		{Key: "lastName", Value: user.GetLastName()},
		{Key: "email", Value: user.GetEmail()}})

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (dbConnector *DbConnector) UpdateUser(userId, newEmail string) (int64, error) {
	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return 0, err
	}

	result, err := dbConnector.UserCollection.UpdateOne(
		dbConnector.Context,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{{"email", newEmail}}},
		},
	)
	if err != nil {
		fmt.Printf("Failed to update user with id %s", userId)
		return 0, err
	}

	return result.ModifiedCount, nil
}

func (dbConnector *DbConnector) InsertEmployee(employee *proto_files.Employee) (string, error) {
	objectId, err := primitive.ObjectIDFromHex(employee.GetUserId())
	if err != nil {
		return "", err
	}

	result, err := dbConnector.EmployeeCollection.InsertOne(dbConnector.Context, bson.D{
		{Key: "userId", Value: objectId},
		{Key: "designation", Value: employee.GetDesignation()}})

	return result.InsertedID.(primitive.ObjectID).Hex(), err
}

func (dbConnector *DbConnector) Close() error {
	fmt.Println("Cleaning up the resources")
	if err := dbConnector.client.Disconnect(dbConnector.Context); err != nil {
		fmt.Println("Failed to disconnect from database")
		return err
	}
	fmt.Println("Sucessfully disconnected from database")
	return nil
}
