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
	apiConfig.CallMapUrl(true)
	assertStrings(apiConfig.next, startingURL, t)
}

func TestCallNextURL(t *testing.T) {
	apiConfig := GenNewApiConfig()

	next, prev, results, err := apiConfig.CallMapUrl(true)
	// call the API once, get the next 20 records
	assertNoError(err, t)
	assertStrings(apiConfig.next, urlPage2, t)
	assertStrings(next, apiConfig.next, t)
	//in the config, next should now be pointing to page 2
	assertStrings(apiConfig.prev, blankURL, t)
	assertStrings(prev, apiConfig.prev, t)
	//in the config, prev should still in blank
	assertVals(len(results), 20, t)

	next2, prev2, results2, err2 := apiConfig.CallMapUrl(true)
	assertNoError(err2, t)
	assertStrings(apiConfig.next, urlPage3, t)
	assertStrings(next2, apiConfig.next, t)
	//in the config, next should now be pointing to page 3
	assertStrings(apiConfig.prev, urlPage1, t)
	assertStrings(prev2, apiConfig.prev, t)
	//in the config, prev should now be pointing to the starting url
	assertVals(len(results2), 20, t)
}

func TestCallPrevUrl(t *testing.T) {
	apiConfig := GenNewApiConfig()

	//assert that trying to get prev url when it doesn't exist in the apiconfig yet should return an error
	_, _, _, err := apiConfig.CallMapUrl(false)
	errorHelpers.AssertError(err, t)

	// call next on the first url, no error should be returned
	firstNext, firstPrev, firstResults, firstCallErr := apiConfig.CallMapUrl(true)
	assertNoError(firstCallErr, t)

	//at this point, prev should be null, as we are at the first page. should return error
	_, _, _, errFromPrev := apiConfig.CallMapUrl(false)
	errorHelpers.AssertError(errFromPrev, t)

	//call the API one more time to get to the second page
	apiConfig.CallMapUrl(true)

	secondNext, secondPrev, secondResults, secondCallErr := apiConfig.CallMapUrl(false)
	assertNoError(secondCallErr, t)
	assertStrings(firstNext, secondNext, t)
	assertStrings(firstPrev, secondPrev, t)
	errorHelpers.AssertReflectDeepEqual(firstResults, secondResults, t)

}

func TestGetPokemonInfo(t *testing.T) {
	type testStruct struct {
		name               string
		url                string
		expectedNumResults int
	}

	pokemonInfoTests := []testStruct{
		{
			"should pass and get 151 records",
			"https://pokeapi.co/api/v2/pokedex/5/",
			151,
		},
		{
			"should fail and get 0 records",
			"https://you failed",
			0,
		},
	}

	for _, test := range pokemonInfoTests {
		t.Run(test.name, func(t *testing.T) {
			got := len(getPokemonInfo(test.url))
			want := test.expectedNumResults
			assertVals(got, want, t)
		})
	}
}
