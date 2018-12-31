package delivery

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Package struct {
	Path     string
	Doorstep string
}


// TODO: should be set in env vars
var supportedDeliveryMethods = [3]string{"dropbox", "email", "requestor"}
var protocol = "http"
var deliveryHostname = "10.0.1.48"
var deliveryPort = "9000"
var deliveryTargetPath = "scan"

var address = protocol + "://" + deliveryHostname + ":" + deliveryPort + "/" + deliveryTargetPath

func Initiate(filepath string, doorstep string) (string, error) {
	jsonData := map[string]string{"path": filepath, "doorstep": doorstep}
	jsonValue, err := json.Marshal(jsonData)
	_, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
	//response, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	return "COLIN", nil
}

func ValidateDestination(destination string) error {
	supportedMethod := false
	for i := range supportedDeliveryMethods {
		if destination == supportedDeliveryMethods[i] {
			supportedMethod = true
			break
		}
	}
	if !supportedMethod {
		return fmt.Errorf("you have requested an unsupported delivery method: %s; please choose from the following: %s", destination, supportedDeliveryMethods )
	}
	return nil
}