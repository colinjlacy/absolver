package delivery

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Package struct {
	Foldername   string `json: foldername`
	EmailAddress string `json: emailAddress`
}

// TODO: should be set in env vars
const protocol = "http"
const deliveryHostname = "10.0.1.48"
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
