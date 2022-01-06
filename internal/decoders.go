package internal

import (
	b64 "encoding/base64"
	"encoding/json"
	"net/http"
)

type protoBody struct {
	Proto_Body string
}

func DecodeJson(bindedStructure interface{}, request *http.Request) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(bindedStructure)
	return err
}

func DecodeProtoBody(request *http.Request) ([]byte, error) {
	protoBody := &protoBody{}
	if err := DecodeJson(protoBody, request); err != nil {
		return []byte{}, err
	}

	return b64.StdEncoding.DecodeString(protoBody.Proto_Body)
}
