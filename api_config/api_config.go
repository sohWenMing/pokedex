package apiconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type apiConfig struct {
	next string
	prev string
}

type JSONResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

const startingURL = "https://pokeapi.co/api/v2/location/"

var blankJSONResponse = JSONResponse{}

func GenNewApiConfig() *apiConfig {
	apiConfig := apiConfig{
		next: startingURL,
		prev: "",
	}
	return &apiConfig
}

func (a *apiConfig) SetConfig(next, prev string) {
	a.next = next
	a.prev = prev
}

func (a *apiConfig) GetNext() (response JSONResponse, err error) {
	if a.next == "" {
		a.resetConfig()
		return blankJSONResponse, errors.New("no further values in map, resetting api destinations")
	}
	res, callNextURLErr := a.callNextURL()
	if callNextURLErr != nil {
		return blankJSONResponse, callNextURLErr
	}
	a.next = res.Next
	a.prev = res.Previous
	return res, nil

}

func (a *apiConfig) resetConfig() {
	a.next = startingURL
	a.prev = ""
}

func (a *apiConfig) callNextURL() (response JSONResponse, err error) {
	if a.next == "" {
		a.resetConfig()
		return blankJSONResponse, errors.New("error: getNext called when there was no next in config")
	}
	res, err := http.Get(a.next)
	checkErr := checkResponseErrAndStatus(res, err)
	if checkErr != nil {
		return blankJSONResponse, checkErr
	}
	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return blankJSONResponse, readErr
	}
	var jsonResponse JSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &jsonResponse)
	if jsonErr != nil {
		return blankJSONResponse, jsonErr
	}
	return jsonResponse, nil
}

func checkResponseErrAndStatus(res *http.Response, errorFromGet error) (err error) {
	if errorFromGet != nil {
		return err
	}
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return fmt.Errorf("error status code: %d status: %s", res.StatusCode, res.Status)
	}
	return nil
}
