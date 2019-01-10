package archive

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

// TODO: should be set in env vars
const protocol = "http"
const archiveHostname = "localhost"
const archivePort = "4000"
const archiveJobsPath = "jobs"
const archiveJobPath = "job"
const archiveImagePath = "image"

const domain = protocol + "://" + archiveHostname + ":" + archivePort + "/"

func PullFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["jobName"]
	imageName := vars["fileName"]
	if jobName == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": "you need to include a job name as a first parameter in your request"}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if imageName == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": "you need to include a file name as a second parameter in your request"}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	address := domain + archiveImagePath
	response, err := http.Get(address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if response.StatusCode > 399 {
		var data map[string]string
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			data = map[string]string{"error": err.Error()}
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(data)
		return
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bodyBytes)
}

func FetchCatalog(w http.ResponseWriter, r *http.Request) {
	address := domain + archiveJobsPath
	response, err := http.Get(address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if response.StatusCode > 399 {
		var data map[string]string
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			data = map[string]string{"error": err.Error()}
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(data)
		return
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bodyBytes)
}

func PullFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobName := vars["jobName"]
	if jobName == "" {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": "you need to include a job name as a parameter in your request"}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	address := domain + archiveJobPath
	response, err := http.Get(address)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		jsonData := map[string]string{"error": err.Error()}
		_ = json.NewEncoder(w).Encode(jsonData)
		return
	}
	if response.StatusCode > 399 {
		var data map[string]string
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			data = map[string]string{"error": err.Error()}
		}
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(data)
		return
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(bodyBytes)
}
