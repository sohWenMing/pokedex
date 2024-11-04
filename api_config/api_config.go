package apiconfig

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	cache "github.com/sohWenMing/pokedex/cache"
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

const startingURL = "https://pokeapi.co/api/v2/location/"

var blankInfo = []string{}

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

func (a *ApiConfig) GetNext(c *cache.Cache) (info []string, err error) {
	next, prev, info, callNextErr := a.callNextURL(c)
	if callNextErr != nil {
		return blankInfo, callNextErr
	}
	a.next = next
	a.prev = prev
	return info, nil

}

func (a *ApiConfig) resetConfig() {
	a.next = startingURL
	a.prev = ""
}

func (a *ApiConfig) callNextURL() (next, prev string, info []string, err error) {

	res, err := http.Get(a.next)
	checkErr := checkResponseErrAndStatus(res, err)
	if checkErr != nil {
		return "", "", blankInfo, checkErr
	}
	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return "", "", blankInfo, readErr
	}
	var jsonResponse JSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &jsonResponse)
	if jsonErr != nil {
		return "", "", blankInfo, jsonErr
	}
	return jsonResponse.Next, jsonResponse.Previous, jsonResponse.Results, nil
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
