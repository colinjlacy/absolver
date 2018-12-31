package main

import (
	"absolver/delivery"
	"absolver/request"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/scan", requestScan).Methods("POST")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func requestScan(w http.ResponseWriter, r *http.Request) {
	// parse req body
	var params request.ScanRequestParams
	_ = json.NewDecoder(r.Body).Decode(&params)
	// validations
	if err := validateRequestParams(params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if err := delivery.ValidateDestination(params.Doorstep); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// send scan request
	fullFilePath, err := request.Attempt(params.Foldername, params.Filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// send delivery request
	deliveryResult, err := delivery.Initiate(fullFilePath, params.Doorstep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// respond
	jsonData := map[string]string{"result": deliveryResult}
	json.NewEncoder(w).Encode(jsonData)
}

func validateRequestParams(params request.ScanRequestParams) error {
	var errors []string
	if params.Doorstep == "" {
		errors = append(errors, `request missing string doorstep property;`)
	}
	if params.Filename == "" {
		errors = append(errors, `request missing string filename property;`)
	}
	if params.Foldername == "" {
		errors = append(errors, `request missing string foldername property;`)
	}
	if len(errors) > 0 {
		return fmt.Errorf("encountered the following validation errors: %s", errors)
	}
	return nil
}
