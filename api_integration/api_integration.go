package apiIntegration

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next string
	Prev string
}

func GetLocation(url string, config *Config) (response LocationResponse, err error) {
	var locationResponse LocationResponse
	res, err := http.Get(url)
	if err != nil {
		return LocationResponse{}, err
	}
	statusError := checkStatus(res)
	if statusError != nil {
		return LocationResponse{}, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	jsonDecodeErr := decoder.Decode(&locationResponse)
	if jsonDecodeErr != nil {
		return LocationResponse{}, jsonDecodeErr
	}
	config.Next = locationResponse.Next
	config.Prev = locationResponse.Previous

	return locationResponse, nil
}

func checkStatus(response *http.Response) (err error) {
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return fmt.Errorf("bad response: %s", response.Status)
	}
	return nil
}
