package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PincodeResponse struct {
	Message    string       `json:"Message"`
	Status     string       `json:"Status"`
	PostOffice []PostOffice `json:"PostOffice"`
}

type PostOffice struct {
	Name     string `json:"Name"`
	State    string `json:"State"`
	District string `json:"District"`
}

func ValidatePincodeandLocation(pincode string, state string, city string) error {

	url := fmt.Sprintf("https://api.postalpincode.in/pincode/%s", pincode)

	response, err := http.Get(url)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	var result []PincodeResponse

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return err
	}

	if len(result) == 0 || result[0].Status != "Success" || len(result[0].PostOffice) == 0 {
		return errors.New("invalid pincode")
	}

	matched := false

	for _, po := range result[0].PostOffice {
		if strings.EqualFold(po.State, state) && strings.EqualFold(po.District, city) {
			matched = true
			break
		}
	}

	if !matched {
		return errors.New("city/state does not match pincode")
	}

	return nil

}
