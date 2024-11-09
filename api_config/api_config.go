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

type MapValue struct {
	Name string
	URL  string
}

const baseURL = "https://pokeapi.co/api/v2/"
const startingURL = "https://pokeapi.co/api/v2/location/"

var blankJsonResults = []MapValue{}

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

func (a *ApiConfig) resetConfig() {
	a.SetConfig(startingURL, "")
}

func (a *ApiConfig) GetNext() string {
	return a.next
}

func (a *ApiConfig) GetPrev() string {
	return a.prev
}

func (a *ApiConfig) SetNext(next string) {
	a.next = next
}

func (a *ApiConfig) SetPrev(prev string) {
	a.prev = prev
}

func (a *ApiConfig) CallMapURL(isNext bool) (next, prev string, results []MapValue, err error) {

	if isNext {
		if a.next == "" {
			a.resetConfig()
			return "", "", blankJsonResults, errors.New("no more locations to show ... resetting")
		}
	}

	if !isNext {
		if a.prev == "" {
			return "", "", blankJsonResults, errors.New("no previous locations to show")
		}
	}

	var urlToCall string

	switch isNext {
	case true:
		urlToCall = a.GetNext()
	case false:
		urlToCall = a.GetPrev()
	}

	res, err := http.Get(urlToCall)
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

	jsonResults := []MapValue{}
	for _, result := range jsonResponse.Results {
		jsonResult := MapValue{
			Name: result.Name,
			URL:  result.URL,
		}
		jsonResults = append(jsonResults, jsonResult)
	}

	a.SetConfig(jsonResponse.Next, jsonResponse.Previous)

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
