package apiconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type ApiConfig struct {
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

type JSONResult struct {
	Name string
	URL  string
}

const startingURL = "https://pokeapi.co/api/v2/location/"

var blankJsonResults = []JSONResult{}

func GenNewApiConfig() *ApiConfig {
	apiConfig := ApiConfig{
		next: startingURL,
		prev: "",
	}
	return &apiConfig
}

func (a *ApiConfig) SetConfig(next, prev string) {
	a.next = next
	a.prev = prev
}

func (a *ApiConfig) GetNext() (next, prev string, results []JSONResult, err error) {
	next, prev, results, callNextErr := a.callNextURL()
	if callNextErr != nil {
		return "", "", []JSONResult{}, callNextErr
	}
	a.next = next
	a.prev = prev
	return next, prev, results, nil

}

func (a *ApiConfig) resetConfig() {
	a.SetConfig(startingURL, "")
}

func (a *ApiConfig) callNextURL() (next, prev string, results []JSONResult, err error) {
	if a.next == "" {
		a.resetConfig()
		return "", "", blankJsonResults, errors.New("no more locations to show ... resetting")
	}

	res, err := http.Get(a.next)
	checkErr := checkResponseErrAndStatus(res, err)
	if checkErr != nil {
		return "", "", blankJsonResults, checkErr
	}
	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", "", blankJsonResults, readErr
	}
	var jsonResponse JSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &jsonResponse)
	if jsonErr != nil {
		return "", "", blankJsonResults, jsonErr
	}
	jsonResults := []JSONResult{}
	for _, result := range jsonResponse.Results {
		jsonResult := JSONResult{
			Name: result.Name,
			URL:  result.URL,
		}
		jsonResults = append(jsonResults, jsonResult)
	}
	return jsonResponse.Next, jsonResponse.Previous, jsonResults, nil
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
