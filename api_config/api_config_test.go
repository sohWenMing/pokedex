package apiconfig

import (
	"testing"

	errorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertVals = errorHelpers.AssertVals
var assertStrings = errorHelpers.AssertStrings
var assertNoError = errorHelpers.AssertNoError

const urlPage1 = "https://pokeapi.co/api/v2/location/?offset=0&limit=20"
const urlPage2 = "https://pokeapi.co/api/v2/location/?offset=20&limit=20"
const urlPage3 = "https://pokeapi.co/api/v2/location/?offset=40&limit=20"
const blankURL = ""

func TestGetNewApiConfig(t *testing.T) {
	apiConfig := GenNewApiConfig()
	assertStrings(apiConfig.next, startingURL, t)
	assertStrings(apiConfig.prev, "", t)

}
func TestResetAPIConfigOnNext(t *testing.T) {
	apiConfig := GenNewApiConfig()
	apiConfig.next = ""
	assertStrings(apiConfig.next, "", t)
	apiConfig.GetNext()
	assertStrings(apiConfig.next, startingURL, t)
}

func TestGetNext(t *testing.T) {
	apiConfig := GenNewApiConfig()

	res, err := apiConfig.GetNext()
	// call the API once, get the next 20 records
	assertNoError(err, t)
	assertStrings(apiConfig.next, urlPage2, t)
	//in the config, next should now be pointing to page 2
	assertStrings(apiConfig.prev, blankURL, t)
	//in the config, prev should still in blank
	assertVals(len(res.Results), 20, t)

	res2, err2 := apiConfig.GetNext()
	assertNoError(err2, t)
	assertStrings(apiConfig.next, urlPage3, t)
	//in the config, next should now be pointing to page 3
	assertStrings(apiConfig.prev, urlPage1, t)
	//in the config, prev should now be pointing to the starting url
	assertVals(len(res2.Results), 20, t)
}
