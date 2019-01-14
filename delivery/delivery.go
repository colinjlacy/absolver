package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EmailPackage struct {
	Foldername   string `json: foldername`
	EmailAddress string `json: emailAddress`
}

type WarehousePackage struct {
	Foldername   string `json: foldername`
	Destination string `json: destination`
}

// TODO: should be set in env vars
const protocol = "http"
const deliveryHostname = "localhost"
const deliveryPort = "9000"
const deliveryTargetPath = "email"

var address = protocol + "://" + deliveryHostname + ":" + deliveryPort + "/" + deliveryTargetPath

func Deliver(foldername string, emailAddress string) error {
	jsonData := map[string]string{"foldername": foldername, "emailAddress": emailAddress}
	jsonValue, err := json.Marshal(jsonData)
	_, err = http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	return nil
}

func Store(foldername string, destination string) error {
	jsonData := map[string]string{"foldername": foldername, "destination": destination}
	jsonValue, err := json.Marshal(jsonData)
	res, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}
	if res.StatusCode > 399 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("could not read byteVessel error response: %s", err)
		}
		return fmt.Errorf("error posting to byteVessel: %s", body)
	}
	return nil
}

