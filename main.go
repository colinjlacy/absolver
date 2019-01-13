package main

import (
	"absolver/archive"
	"absolver/delivery"
	"absolver/request"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type DeliveryInstructions struct {
	Foldername   string `json: foldername`
	EmailAddress string `json: emailAddress`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/scan", requestScan).Methods("POST")
	router.HandleFunc("/email", emailDelivery).Methods("POST")
	router.HandleFunc("/jobs", archive.FetchCatalog).Methods("GET")
	router.HandleFunc("/job/{jobName}", archive.PullFolder).Methods("GET")
	router.HandleFunc("/job/{jobName}", archive.DeleteFolder).Methods("DELETE")
	router.HandleFunc("/image/{jobName}/{fileName}", archive.PullFile).Methods("GET")
	router.HandleFunc("/image/{jobName}/{fileName}", archive.RemoveFile).Methods("DELETE")

	headersOk := handlers.AllowedHeaders([]string{"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*", "localhost", "localhost:4200"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "OPTIONS"})
	corsRouter := handlers.CORS(headersOk, originsOk, methodsOk)(router)

	log.Fatal(http.ListenAndServe(":3000", corsRouter))
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
	// send scan request
	response, err := request.Attempt(params.Foldername, params.Filename, params.PrettyName)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// respond
	jsonData := map[string]string{"filename": response.Filename, "thumbnail": response.Thumbnail}
	json.NewEncoder(w).Encode(jsonData)
}

func emailDelivery(w http.ResponseWriter, r *http.Request) {
	// parse req body
	var params DeliveryInstructions
	_ = json.NewDecoder(r.Body).Decode(&params)
	// TODO: validate request body!!!
	// send delivery request
	err := delivery.Deliver(params.Foldername, params.EmailAddress)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(params)
}

func validateRequestParams(params request.ScanRequestParams) error {
	var errors []string
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
