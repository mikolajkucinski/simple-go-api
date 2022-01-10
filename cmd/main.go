package main

import (
	"awesomeProject/internal"
	proto_files "awesomeProject/internal/proto-files"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

var dbConnector *internal.DbConnector

func getHandler(responseWriter http.ResponseWriter, request *http.Request) {
	protoBody, err := ioutil.ReadAll(request.Body)

	getRequest := &proto_files.GetRequest{}
	if err := proto.Unmarshal(protoBody, getRequest); err != nil {
		log.Fatalf("Failed to unmarshall the GetRequest")
		return
	}

	objId, err := primitive.ObjectIDFromHex(getRequest.GetUserId())
	if err != nil {
		log.Fatalf("Failed to parse string to ObjectId")
		return
	}

	employee, err := dbConnector.FindEmployeeByUserId(objId)
	if err != nil {
		log.Fatalf("Failed to retrieve employee from the database")
		return
	}
	user, err := dbConnector.FindUserById(objId)
	if err != nil {
		log.Fatalf("Failed to retrieve user from the database")
		return
	}

	getResponse := &proto_files.GetResponse{
		Firstname:   user.FirstName,
		Lastname:    user.LastName,
		Email:       user.Email,
		EmployeeId:  employee.Id.Hex(),
		Designation: employee.Designation,
	}
	getResponseMarshalled, err := proto.Marshal(getResponse)
	if err != nil {
		log.Fatalf("Failed to marshal the GetResponse")
		return
	}
	responseWriter.Write(getResponseMarshalled)
}

func postHandler(responseWriter http.ResponseWriter, request *http.Request) {
	protoBody, err := ioutil.ReadAll(request.Body)

	postRequest := &proto_files.PostRequest{}
	if err := proto.Unmarshal(protoBody, postRequest); err != nil {
		log.Fatalf("Failed to unmarshall the PostRequest")
		return
	}

	userId, err := dbConnector.InsertUser(postRequest.GetFirstName(), postRequest.GetLastName(), postRequest.GetEmail())
	if err != nil {
		log.Fatalf("Failed to insert user into database, reason: %s", err.Error())
		return
	}
	_, err = dbConnector.InsertEmployee(userId, postRequest.GetDesignation())
	if err != nil {
		log.Fatalf("Failed to insert employee into database")
		return
	}

	postResponse := &proto_files.PostResponse{Id: userId.Hex()}
	getResponseMarshalled, err := proto.Marshal(postResponse)
	if err != nil {
		log.Fatalf("Failed to marshal the PostResponse")
		return
	}
	responseWriter.Write(getResponseMarshalled)
}

func patchHandler(responseWriter http.ResponseWriter, request *http.Request) {
	protoBody, err := ioutil.ReadAll(request.Body)

	patchRequest := &proto_files.PatchRequest{}
	if err := proto.Unmarshal(protoBody, patchRequest); err != nil {
		log.Fatalf("Failed to unmarshall the PostRequest")
		return
	}

	objId, err := primitive.ObjectIDFromHex(patchRequest.GetId())
	if err != nil {
		log.Fatalf("Failed to parse string to ObjectId")
		return
	}
	_, err = dbConnector.UpdateUser(objId, patchRequest.GetEmail())
	if err != nil {
		log.Fatalf("Failed to update user")
		return
	}

	responseWriter.Write([]byte("Sucessfully updated user!\n"))
}

func main() {
	dbConnector = &internal.DbConnector{}
	dbConnector.Connect()
	defer dbConnector.Close()

	r := mux.NewRouter()
	r.HandleFunc("/assignment/user", getHandler).Methods("GET")
	r.HandleFunc("/assignment/user", postHandler).Methods("POST")
	r.HandleFunc("/assignment/user", patchHandler).Methods("PATCH")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8000", r))
}
