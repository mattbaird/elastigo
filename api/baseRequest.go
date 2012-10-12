package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func DoCommand(method string, url string, data interface{}) (string, error) {
	var response map[string]interface{}
	var body string
	req, err := ElasticSearchRequest(method, url)
	if err != nil {
		return body, err
	}
	if data != nil {
		err = req.SetBodyJson(data)
		if err != nil {
			return body, err
		}
	}
	body, err = req.Do(&response)
	if err != nil {
		return body, err
	}
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		return body, errors.New(fmt.Sprintf("Error [%s] Status [%s]", error, status))
	}
	return body, nil
}

// The API also allows to check for the existance of a document using HEAD
// This appears to be broken in the current version of elasticsearch 0.19.10, currently
// returning nothing
func Exists(pretty bool, index string, _type string, id string) (BaseResponse, error) {
	var response map[string]interface{}
	var body string
	var url string
	var retval BaseResponse

	if len(_type) > 0 {
		url = fmt.Sprintf("/%s/%s/%s?%s", index, _type, id, Pretty(pretty))
	} else {
		url = fmt.Sprintf("/%s/%s?%s", index, id, Pretty(pretty))
	}
	req, err := ElasticSearchRequest("HEAD", url)
	if err != nil {
		// some sort of generic error handler		
	}
	body, err = req.Do(&response)
	if error, ok := response["error"]; ok {
		status, _ := response["status"]
		log.Fatalf("Error: %v (%v)\n", error, status)
	} else {
		// marshall into json
		jsonErr := json.Unmarshal([]byte(body), &retval)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}
	}
	fmt.Println(body)
	return retval, err
}
