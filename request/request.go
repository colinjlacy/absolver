package request

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ScanRequestParams struct {
	Filename         string `json: filename`
	Foldername       string `json: foldername`
	IncludeThumbnail bool   `json: includeThumbnail`
}

type ScanRequestResponse struct {
	Filename  string `json: filename`
	Thumbnail string `json: thumbnail`
}

// TODO: should be set in env vars
var protocol = "http"
var scantasticHostname = "10.0.1.48"
var scantasticPort = "8000"
var scantasticTargetPath = "scan"

var address = protocol + "://" + scantasticHostname + ":" + scantasticPort + "/" + scantasticTargetPath
var busy = false

func Attempt(foldername string, filename string) (ScanRequestResponse, error) {
	// TODO: customize error types
	if busy {
		return ScanRequestResponse{}, fmt.Errorf("busy")
	}
	busy = true
	jsonData := ScanRequestParams{Filename:filename, Foldername:foldername, IncludeThumbnail:true}
	jsonValue, err := json.Marshal(jsonData)
	if err != nil {
		return ScanRequestResponse{}, fmt.Errorf("could not marshal JSON from request attempt parameters: %s", err)
	}
	response, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	busy = false
	if err != nil {
		return ScanRequestResponse{}, fmt.Errorf("there was a problem with the request to the scanner controller: %s", err)
	}
	var data ScanRequestResponse
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		return ScanRequestResponse{}, fmt.Errorf("could not read request response body: %s", err)
	}
	// TODO: this isn't returning response error data correctly
	if response.StatusCode > 399 {
		return ScanRequestResponse{}, fmt.Errorf("the request to the scanner controller reruend a status of %d: %s", response.StatusCode, data)
	}
	return data, nil
}
