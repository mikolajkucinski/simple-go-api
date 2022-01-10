package main

import (
	proto_files "awesomeProject/internal/proto-files"
	"bytes"
	"fmt"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"log"
	"net/http"
)

func makePost(request *proto_files.PostRequest) *proto_files.PostResponse {
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}

	resp, err := http.Post("http://127.0.0.1:8000/assignment/user", "application/x-binary", bytes.NewReader(req))
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	respObj := &proto_files.PostResponse{}
	proto.Unmarshal(respBytes, respObj)
	return respObj
}

func makeGet(request *proto_files.GetRequest) *proto_files.GetResponse {
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}
	httpReq, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8000/assignment/user", bytes.NewReader(req))
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Unable to read bytes from request : %v", err)
	}

	respObj := &proto_files.GetResponse{}
	proto.Unmarshal(respBytes, respObj)
	return respObj
}

func makePatch(request *proto_files.PatchRequest) {
	req, err := proto.Marshal(request)
	if err != nil {
		log.Fatalf("Unable to marshal request : %v", err)
	}
	httpReq, err := http.NewRequest(http.MethodPatch, "http://127.0.0.1:8000/assignment/user", bytes.NewReader(req))
	if err != nil {
		log.Fatal(err)
	}
	_, err = http.DefaultClient.Do(httpReq)
	if err != nil {
		log.Fatalf("Unable to read from the server : %v", err)
	}
}

func main() {

	postRequest := &proto_files.PostRequest{
		FirstName:   "Jan",
		LastName:    "Kowalski",
		Email:       "jan.kowalski@nokia.com,",
		Designation: "Node.js developer",
	}
	postResponse := makePost(postRequest)
	fmt.Printf("UserId from the response = %s\n", postResponse.GetId())

	getRequest := &proto_files.GetRequest{UserId: postResponse.GetId()}
	getResponse := makeGet(getRequest)
	fmt.Println(getResponse)

	patchRequest := &proto_files.PatchRequest{Id: postResponse.GetId(), Email: "none"}
	makePatch(patchRequest)

	getRequest = &proto_files.GetRequest{UserId: postResponse.GetId()}
	getResponse = makeGet(getRequest)
	fmt.Println(getResponse)
}
