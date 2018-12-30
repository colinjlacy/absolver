package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type RequestParams struct {
	Filename   string `json: filename`
	Foldername string `json: foldername`
	Doorstep   string `json: doorstep`
}

type RequestResponse struct {
	Path string `json: path`
}

// TODO: should be set in env vars
var protocol = "http"
var scantasticHostname = "10.0.1.48"
var scantasticPort = "8000"
var scantasticTargetPath = "scan"

var address = protocol + "://" + scantasticHostname + ":" + scantasticPort + "/" + scantasticTargetPath
var busy = false

func Attempt(filepath string, filename string) (string, error) {
	// TODO: customize error types
	if busy {
		return "", fmt.Errorf("busy")
	}
	busy = true
	jsonData := map[string]string{"foldername": filepath, "filename": filename}
	log.Println(jsonData)
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return "", fmt.Errorf("could not marshal JSON from request attempt parameters: %s", err)
	}
	response, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	log.Println(response.StatusCode)
	busy = false
	if err != nil {
		return "", fmt.Errorf("there was a problem with the request to the scanner controller: %s", err)
	}
	if response.StatusCode > 399 {
		return "", fmt.Errorf("the request to the scanner controller reruend a status of %d: %s", response.StatusCode, err)
	}
	var data RequestResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	log.Println("scanner response data", data)
	if err != nil {
		return "", fmt.Errorf("could not read request response body: %s", err)
	}
	return data.Path, nil
}
