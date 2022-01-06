package internal

import (
	proto_files "awesomeProject/internal/proto-files"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net/http"
)

func DecodeEmployee(protoBody string) (*proto_files.Employee, error) {
	decodedProtoBody, _ := b64.StdEncoding.DecodeString(protoBody)
	employee := &proto_files.Employee{}
	err := proto.Unmarshal(decodedProtoBody, employee)

	return employee, err
}

func DecodeJson(bindedStructure interface{}, request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(bindedStructure)
	if err != nil {
		fmt.Println("Failed to decode json")
	}

	return err
}
