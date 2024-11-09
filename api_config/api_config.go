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

type MapJsonResponse struct {
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

type ExploreJSONResponse struct {
	Areas []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"areas"`
	GameIndices []struct {
		GameIndex  int `json:"game_index"`
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
	} `json:"game_indices"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
}

type RegionJSONResponse struct {
	ID        int `json:"id"`
	Locations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"locations"`
	MainGeneration struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"main_generation"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	Pokedexes []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"pokedexes"`
	VersionGroups []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"version_groups"`
}

type PokedexJsonResponse struct {
	Descriptions []struct {
		Description string `json:"description"`
		Language    struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"descriptions"`
	ID           int    `json:"id"`
	IsMainSeries bool   `json:"is_main_series"`
	Name         string `json:"name"`
	Names        []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEntries []struct {
		EntryNumber    int `json:"entry_number"`
		PokemonSpecies struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon_species"`
	} `json:"pokemon_entries"`
	Region struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"region"`
	VersionGroups []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"version_groups"`
}

type PokemonEntry struct {
	EntryNumber    int
	PokemonSpecies struct {
		Name string
		URL  string
	}
}

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

func (a *ApiConfig) CallMapUrl(isNext bool) (next, prev string, results []MapValue, err error) {

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

	bodyBytes, err := getBodyBytes(urlToCall)
	if err != nil {
		return "", "", blankJsonResults, err
	}
	var jsonResponse MapJsonResponse
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

func GetPokemonByLocation(location string) (pokedexURLS []string, err error) {
	regionUrl, regionErr := callExploreURL(location)
	if regionErr != nil {
		return []string{}, regionErr
	}

	pokedexUrls, pokedexUrlsErr := callRegionURLToGetPokedexURLS(regionUrl)

	if pokedexUrlsErr != nil {
		return []string{}, pokedexUrlsErr
	}

	return pokedexUrls, nil

}

func getPokemonInfo(url string) (pokemonEntries []PokemonEntry) {
	emptyEntries := []PokemonEntry{}
	bodyBytes, bodyBytesErr := getBodyBytes(url)
	if bodyBytesErr != nil {
		return emptyEntries
	}
	var pokedexJsonResponse PokedexJsonResponse
	jsonErr := json.Unmarshal(bodyBytes, &pokedexJsonResponse)
	if jsonErr != nil {
		return emptyEntries
	}
	returnedEntries := []PokemonEntry{}
	for _, entry := range pokedexJsonResponse.PokemonEntries {
		pokemonEntry := PokemonEntry{
			EntryNumber: entry.EntryNumber,
			PokemonSpecies: struct {
				Name string
				URL  string
			}{
				entry.PokemonSpecies.Name, entry.PokemonSpecies.URL,
			},
		}
		returnedEntries = append(returnedEntries, pokemonEntry)
	}
	return returnedEntries
}

func callExploreURL(location string) (regionUrl string, err error) {
	urlToCall := startingURL + location

	bodyBytes, bodyBytesErr := getBodyBytes(urlToCall)
	if bodyBytesErr != nil {
		return "", bodyBytesErr
	}
	var exploreJsonResponse ExploreJSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &exploreJsonResponse)
	if jsonErr != nil {
		return "", jsonErr
	}
	return exploreJsonResponse.Region.URL, nil
}

func callRegionURLToGetPokedexURLS(urlToCall string) (pokedexURLS []string, err error) {
	bodyBytes, bodyBytesErr := getBodyBytes(urlToCall)
	if bodyBytesErr != nil {
		return []string{}, bodyBytesErr
	}
	var regionJSONResponse RegionJSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &regionJSONResponse)
	if jsonErr != nil {
		return []string{}, jsonErr
	}
	returnedUrls := []string{}
	for _, pokedex := range regionJSONResponse.Pokedexes {
		returnedUrls = append(returnedUrls, pokedex.URL)
	}
	return returnedUrls, nil
}

func getBodyBytes(url string) (bodyBytes []byte, err error) {
	res, responseErr := http.Get(url)
	if responseErr != nil {
		return []byte{}, responseErr
	}

	checkErr := checkResponseErrAndStatus(res, responseErr)

	if checkErr != nil {
		return []byte{}, checkErr
	}
	bodyBytes, readErr := io.ReadAll(res.Body)
	if readErr != nil {
		return []byte{}, readErr
	}
	return bodyBytes, nil

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
