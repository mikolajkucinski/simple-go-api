package internal

import (
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

type Employee struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	Designation string             `json:"designation" bson:"designation"`
}

type User struct {
	Id        primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName string             `json:"firstName" bson:"firstName"`
	LastName  string             `json:"lastName" bson:"lastName"`
	Email     string             `json:"email" bson:"email"`
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

func (dbConnector *DbConnector) FindEmployeeByUserId(userId primitive.ObjectID) (Employee, error) {
	employee := &Employee{}
	if err := dbConnector.EmployeeCollection.FindOne(dbConnector.Context, bson.M{"userId": userId}).Decode(employee); err != nil {
		return Employee{}, err
	}

	return *employee, nil
}

func (dbConnector *DbConnector) FindUserById(id primitive.ObjectID) (User, error) {
	user := &User{}
	if err := dbConnector.UserCollection.FindOne(dbConnector.Context, bson.M{"_id": id}).Decode(user); err != nil {
		return User{}, err
	}

	return *user, nil
}

func (dbConnector *DbConnector) InsertUser(firstName, lastName, email string) (primitive.ObjectID, error) {
	result, err := dbConnector.UserCollection.InsertOne(dbConnector.Context, bson.D{
		{Key: "firstName", Value: firstName},
		{Key: "lastName", Value: lastName},
		{Key: "email", Value: email}})
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (dbConnector *DbConnector) InsertEmployee(userId primitive.ObjectID, designation string) (primitive.ObjectID, error) {
	result, err := dbConnector.EmployeeCollection.InsertOne(dbConnector.Context, bson.D{
		{Key: "userId", Value: userId},
		{Key: "designation", Value: designation}})
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return result.InsertedID.(primitive.ObjectID), err
}

func (dbConnector *DbConnector) UpdateUser(userId primitive.ObjectID, email string) (int64, error) {
	result, err := dbConnector.UserCollection.UpdateOne(
		dbConnector.Context,
		bson.M{"_id": userId},
		bson.D{
			{"$set", bson.D{{"email", email}}},
		},
	)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, err
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
