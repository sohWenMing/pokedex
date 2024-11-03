package cache

import (
	"sync"
	"testing"
	"time"

	errorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertVals = errorHelpers.AssertVals
var AssertReflectDeepEqual = errorHelpers.AssertReflectDeepEqual

func TestNewCache(t *testing.T) {
	testDuration := 10 * time.Millisecond
	cache := NewCache(testDuration)
	counterChan := make(chan struct{})
	doneChan := make(chan struct{})

	counterVal := 0
	var wg sync.WaitGroup
	wg.Add(10)
	//create a waitgroup that is waiting for 10 signals

	go cache.ActivateCacheClear(counterChan, doneChan)

	go func() {
		for {
			select {
			case <-doneChan:
				return
			case <-counterChan:
				counterVal++
				wg.Done()
			}
		}
	}()

	timeout := time.After(2 * time.Second)

	wgDoneChan := make(chan struct{})

	go func() {
		wg.Wait()
		wgDoneChan <- struct{}{}
	}()

	select {
	case <-timeout:
		t.Errorf("Timeout occured before 10 counts could finish")
	case <-wgDoneChan:
		assertVals(counterVal, 10, t)
	}

}

func TestWriteToCache(t *testing.T) {

	//not testing for clearing of cache yet, validity doesn't matter
	type keyValue struct {
		key    string
		values []string
	}

	type testStruct struct {
		name               string
		testKeyValues      []keyValue
		expectedUrlsToInfo map[string][]string
	}

	testKeyValues := []keyValue{
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
