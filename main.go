package main

import (
	"absolver/delivery"
	"absolver/request"
	"encoding/json"
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
	var params request.RequestParams
	_ = json.NewDecoder(r.Body).Decode(&params)
	fullFilePath, err := request.Attempt(params.Foldername, params.Filename)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		jsonValue, _ := json.Marshal(jsonData)
		json.NewEncoder(w).Encode(jsonValue)
		return
	}
	deliveryResult, err := delivery.Initiate(fullFilePath, params.Doorstep)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		jsonValue, _ := json.Marshal(jsonData)
		json.NewEncoder(w).Encode(jsonValue)
		return
	}
	jsonData := map[string]string{"result": deliveryResult}
	jsonValue, _ := json.Marshal(jsonData)
	json.NewEncoder(w).Encode(jsonValue)
}
