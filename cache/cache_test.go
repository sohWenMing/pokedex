package cache

import (
	"sync"
	"testing"
	"time"

	errorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var AssertReflectDeepEqual = errorHelpers.AssertReflectDeepEqual

func TestWriteToCache(t *testing.T) {

	//not testing for clearing of cache yet, validity doesn't matter
	type keyValue struct {
		key    string
		values []string
	}

	var testKeyValues = []keyValue{
		{
			key: "url1",
			values: []string{
				"test string 1", "test string 2",
			},
		},
		{
			key: "url2",
			values: []string{
				"test string 3", "test string 4",
			},
		},
	}

	type testStruct struct {
		name               string
		testKeyValues      []keyValue
		expectedUrlsToInfo map[string][]string
	}

	tests := []testStruct{
		{
			name:          "write two unique values - no overwrite",
			testKeyValues: testKeyValues,
			expectedUrlsToInfo: map[string][]string{
				testKeyValues[0].key: testKeyValues[0].values,
				testKeyValues[1].key: testKeyValues[1].values,
			},
		},
		{
			name:          "write same url twice, should overwrite",
			testKeyValues: []keyValue{testKeyValues[0], testKeyValues[0]},
			expectedUrlsToInfo: map[string][]string{
				testKeyValues[0].key: testKeyValues[0].values,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewCache(0 * time.Second)

			for _, keyValue := range test.testKeyValues {
				cache.WriteToCache(keyValue.key, keyValue.values)

			}
			gotUrlsToInfo := make(map[string][]string)
			for key, cacheMapVal := range cache.cacheMap {
				gotUrlsToInfo[key] = cacheMapVal.info
			}
			AssertReflectDeepEqual(gotUrlsToInfo, test.expectedUrlsToInfo, t)
		})
	}

}

func TestWriteAndClearingCache(t *testing.T) {

	type keyValue struct {
		key    string
		values []string
	}

	var testKeyValues = []keyValue{
		{
			key: "url1",
			values: []string{
				"test string 1", "test string 2",
			},
		},
		{
			key: "url2",
			values: []string{
				"test string 3", "test string 4",
			},
		},
	}
	testDuration := 1 * time.Millisecond
	//this should make the cache check every 10 seconds
	cache := NewCache(testDuration)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go cache.ActivateCacheClear(doneChan)

	for _, testKeyValue := range testKeyValues {
		time.Sleep(600 * time.Millisecond)
		cache.WriteToCache(testKeyValue.key, testKeyValue.values)
		wg.Done()
	}
	wg.Wait()
	doneChan <- struct{}{}

	got := []keyValue{}
	for key, cacheMapVal := range cache.cacheMap {
		workingKeyValue := keyValue{
			key:    key,
			values: cacheMapVal.info,
		}
		got = append(got, workingKeyValue)
	}
	expected := []keyValue{}
	expected = append(expected, testKeyValues[1])
	AssertReflectDeepEqual(got, expected, t)
}
