package main

import (
	"awesomeProject/internal"
	proto_files "awesomeProject/internal/proto-files"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var dbConnector *internal.DbConnector

func getHandler(responseWriter http.ResponseWriter, request *http.Request) {
	rawQuery := request.URL.Query()
	protoBody, isPresent := rawQuery["proto_body"]
	if !isPresent {
		fmt.Println("Failed to retrieve proto_body")
		return
	}

	employee, err := internal.DecodeEmployee(protoBody[0])
	if err != nil {
		fmt.Println("Failed to decode employee")
	}

	result, err := dbConnector.FindUserByUserId(employee.UserId)
	if err != nil {
		fmt.Printf("Failed to retrieve data for user id = %s", employee.GetUserId())
		return
	}

	retrievedUser := proto_files.User{
		FirstName: result[0]["firstName"].(string),
		LastName:  result[0]["secondName"].(string),
		Email:     result[0]["email"].(string),
	}
	fmt.Printf("User retrieved = %s", retrievedUser)
}

func postHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decodedJson := &struct {
		FirstName   string
		LastName    string
		Email       string
		Designation string
	}{}
	if err := internal.DecodeJson(decodedJson, request); err != nil {
		return
	}

	insertedUserId, err := dbConnector.InsertUser(&proto_files.User{
		FirstName: decodedJson.FirstName,
		LastName:  decodedJson.LastName,
		Email:     decodedJson.Email})
	if err != nil {
		fmt.Println("Failed to insert new user")
		return
	}

	insertedEmployeeId, err := dbConnector.InsertEmployee(&proto_files.Employee{
		UserId:      insertedUserId,
		Designation: decodedJson.Designation,
	})
	if err != nil {
		fmt.Println("Failed to insert new employee")
		return
	}

	fmt.Printf("Inserted new user and employers under ids %s and %s", insertedUserId, insertedEmployeeId)
}

func patchHandler(responseWriter http.ResponseWriter, request *http.Request) {
	decodedJson := &struct {
		UserId string
		Email  string
	}{}
	if err := internal.DecodeJson(decodedJson, request); err != nil {
		return
	}

	_, err := dbConnector.UpdateUser(decodedJson.UserId, decodedJson.Email)
	if err != nil {
		fmt.Printf("Failed to update user with id = %s", decodedJson.UserId)
		return
	}

	fmt.Printf("Updated user with id = %s", decodedJson.UserId)
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
