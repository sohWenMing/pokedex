package apiconfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

/*
======================================================

# JSON STRUCTS - START

======================================================
*/
type MapJSONResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type LocationRegionJSONResponse struct {
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

type PokedexJSONResponse struct {
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

type ExploreJSONResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int           `json:"chance"`
				ConditionValues []interface{} `json:"condition_values"`
				MaxLevel        int           `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

/*
======================================================

JSON STRUCTS - END

======================================================
*/

/*
======================================================

# OTHER STRUCT DEFINITIONS - START

======================================================
*/
type ApiConfig struct {
	next string
	prev string
}

type MapValue struct {
	Name string
	URL  string
}

type PokemonEntry struct {
	EntryNumber    int
	PokemonSpecies struct {
		Name string
		URL  string
	}
}

type Pokemon struct {
	Name string
	URL  string
}

/*
======================================================

OTHER STRUCT DEFINITIONS - END

======================================================
*/

const baseURL = "https://pokeapi.co/api/v2"
const locationURL = baseURL + "/location/"
const exploreURL = baseURL + "/location-area/"

/*
======================================================

# API CONFIGURATIONS - START

======================================================
*/
func GenNewApiConfig() *ApiConfig {
	apiConfig := ApiConfig{
		next: locationURL,
		prev: "",
	}
	return &apiConfig
}
func (a *ApiConfig) SetConfig(next, prev string) {
	a.next = next
	a.prev = prev
}

func (a *ApiConfig) resetConfig() {
	a.SetConfig(locationURL, "")
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

/*
======================================================

API CONFIGURATIONS - END

======================================================
*/

func CallExploreUrl(area string) (pokemon []Pokemon, err error) {

	emptyPokemon := []Pokemon{}
	urlToCall := exploreURL + area
	bodyBytes, bodyBytesErr := getBodyBytes(urlToCall)
	if bodyBytesErr != nil {
		return emptyPokemon, bodyBytesErr
	}
	var exploreJSONResponse ExploreJSONResponse
	jsonErr := json.Unmarshal(bodyBytes, &exploreJSONResponse)
	if jsonErr != nil {
		return emptyPokemon, jsonErr
	}

	pokemonReturned := []Pokemon{}
	pokemonEncounters := exploreJSONResponse.PokemonEncounters
	for _, pokemon := range pokemonEncounters {
		currentPokemon := Pokemon{
			pokemon.Pokemon.Name,
			pokemon.Pokemon.URL,
		}
		pokemonReturned = append(pokemonReturned, currentPokemon)
	}
	return pokemonReturned, nil
}

func (a *ApiConfig) CallMapUrl(isNext bool) (next, prev string, results []MapValue, err error) {
	blankJsonResults := []MapValue{}
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
	var jsonResponse MapJSONResponse
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

func GetPokemonByLocationRegion(location string) (pokemon []PokemonEntry, err error) {
	regionUrl, regionErr := callLocationRegionURL(location)
	if regionErr != nil {
		return []PokemonEntry{}, regionErr
	}

	pokedexUrls, pokedexUrlsErr := callRegionURLToGetPokedexURLS(regionUrl)

	if pokedexUrlsErr != nil {
		return []PokemonEntry{}, pokedexUrlsErr
	}

	return getPokemonInfoByPokedexUrls(pokedexUrls), nil
}

func getPokemonInfoByPokedexUrls(pokedexUrls []string) (pokemon []PokemonEntry) {

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(pokedexUrls))

	pokemonChan := make(chan []PokemonEntry)

	for _, pokedexUrl := range pokedexUrls {
		url := pokedexUrl
		go func(url string) {
			defer waitGroup.Done()
			pokemonEntries := getPokemonInfo(url)
			pokemonChan <- pokemonEntries
		}(url)
	}

	go func() {
		waitGroup.Wait()
		close(pokemonChan)
	}()

	for pokeData := range pokemonChan {
		pokemon = append(pokemon, pokeData...)
	}

	return pokemon
}

func getPokemonInfo(url string) (pokemonEntries []PokemonEntry) {
	emptyEntries := []PokemonEntry{}
	bodyBytes, bodyBytesErr := getBodyBytes(url)
	if bodyBytesErr != nil {
		return emptyEntries
	}
	var pokedexJsonResponse PokedexJSONResponse
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

func callLocationRegionURL(area string) (regionUrl string, err error) {
	urlToCall := locationURL + area

	bodyBytes, bodyBytesErr := getBodyBytes(urlToCall)
	if bodyBytesErr != nil {
		return "", bodyBytesErr
	}
	var exploreJsonResponse LocationRegionJSONResponse
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
	if res != nil {

		defer res.Body.Close()
	}
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
