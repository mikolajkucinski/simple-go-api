package main

import (
	"awesomeProject/internal"
	proto_files "awesomeProject/internal/proto-files"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

var dbConnector *internal.DbConnector

func getHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decodedProtoBody, err := internal.DecodeProtoBody(request)
	if err != nil {
		fmt.Println("Failed to decode proto body")
		return
	}

	getRequest := &proto_files.GetRequest{}
	if err := proto.Unmarshal(decodedProtoBody, getRequest); err != nil {
		fmt.Println("Failed to unmarshall the GetRequest")
		return
	}

	objId, err := primitive.ObjectIDFromHex(getRequest.GetUserId())
	if err != nil {
		fmt.Println("Failed to parse string to ObjectId")
		return
	}

	employee, err := dbConnector.FindEmployeeByUserId(objId)
	if err != nil {
		fmt.Println("Failed to retrieve employee from the database")
		return
	}
	user, err := dbConnector.FindUserById(objId)
	if err != nil {
		fmt.Println("Failed to retrieve user from the database")
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
		fmt.Println("Failed to marshal the GetResponse")
		return
	}
	responseWriter.Write(getResponseMarshalled)
}

func postHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decodedProtoBody, err := internal.DecodeProtoBody(request)
	if err != nil {
		fmt.Println("Failed to decode proto body")
		return
	}

	postRequest := &proto_files.PostRequest{}
	if err := proto.Unmarshal(decodedProtoBody, postRequest); err != nil {
		fmt.Println("Failed to unmarshall the PostRequest")
		return
	}

	userId, err := dbConnector.InsertUser(postRequest.GetFirstName(), postRequest.GetLastName(), postRequest.GetEmail())
	if err != nil {
		fmt.Printf("Failed to insert user into database, reason: %s", err.Error())
		return
	}
	_, err = dbConnector.InsertEmployee(userId, postRequest.GetDesignation())
	if err != nil {
		fmt.Println("Failed to insert employee into database")
		return
	}

	postResponse := &proto_files.PostResponse{Id: userId.Hex()}
	getResponseMarshalled, err := proto.Marshal(postResponse)
	if err != nil {
		fmt.Println("Failed to marshal the PostResponse")
		return
	}
	responseWriter.Write(getResponseMarshalled)
}

func patchHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decodedProtoBody, err := internal.DecodeProtoBody(request)
	if err != nil {
		fmt.Println("Failed to decode proto body")
		return
	}

	patchRequest := &proto_files.PatchRequest{}
	if err := proto.Unmarshal(decodedProtoBody, patchRequest); err != nil {
		fmt.Println("Failed to unmarshall the PostRequest")
		return
	}

	objId, err := primitive.ObjectIDFromHex(patchRequest.GetId())
	if err != nil {
		fmt.Println("Failed to parse string to ObjectId")
		return
	}
	_, err = dbConnector.UpdateUser(objId, patchRequest.GetEmail())
	if err != nil {
		fmt.Println("Failed to update user")
		return
	}

	responseWriter.Write([]byte("Sucessfully updated user!\n"))
}

func main() {
	//post := &proto_files.PostRequest{
	//	FirstName:   "Izabelka",
	//	LastName:    "Wolek",
	//	Email:       "none",
	//	Designation: "Accountant",
	//}
	//
	//result, _ := proto.Marshal(post)
	//sEnc := b64.StdEncoding.EncodeToString(result)
	//fmt.Println(sEnc)

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
