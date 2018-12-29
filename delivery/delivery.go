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

var supportedDeliveryMethods = [3]string{"dropbox", "email", "requestor"}

// TODO: should be set in env vars
var protocol = "http"
var deliveryHostname = "localhost"
var deliveryPort = "9000"
var deliveryTargetPath = "scan"

var address = protocol + "://" + deliveryHostname + ":" + deliveryPort + "/" + deliveryTargetPath

func Initiate(filepath string, doorstep string) (string, error) {
	err := validateDestination(doorstep)
	if err != nil {
		return "", err
	}
	jsonData := map[string]string{"path": filepath, "doorstep": doorstep}
	jsonValue, err := json.Marshal(jsonData)
	_, err = http.Post("https://httpbin.org/post", "application/json", bytes.NewBuffer(jsonValue))
	//response, err := http.Post(address, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	return "COLIN", nil
}

func validateDestination(destination string) error {
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