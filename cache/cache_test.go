package cache

import (
	"slices"
	"sync"
	"testing"
	"time"

	errorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

type valStruct struct {
	key    string
	next   string
	prev   string
	values []string
}

var testVals = []valStruct{
	{
		key:  "url1",
		next: "next",
		prev: "prev",
		values: []string{
			"test string 1", "test string 2",
		},
	},
	{
		key:  "url2",
		next: "next",
		prev: "prev",
		values: []string{
			"test string 3", "test string 4",
		},
	},
}
var AssertReflectDeepEqual = errorHelpers.AssertReflectDeepEqual

func TestWriteToCache(t *testing.T) {

	//not testing for clearing of cache yet, validity doesn't matter

	type testStruct struct {
		name         string
		testVals     []valStruct
		expectedVals []valStruct
	}

	tests := []testStruct{
		{
			name:         "write two unique values - no overwrite",
			testVals:     testVals,
			expectedVals: testVals,
		},
		{
			name:         "write same url twice, should overwrite",
			testVals:     []valStruct{testVals[0], testVals[0]},
			expectedVals: []valStruct{testVals[0]},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cache := NewCache(0 * time.Second)

			for _, testVal := range test.testVals {
				cache.WriteToCache(testVal.key, testVal.next, testVal.prev, testVal.values)

			}
			gotVals := []valStruct{}
			keys := []string{}
			for key := range cache.cacheMap {
				keys = append(keys, key)
			}
			slices.Sort(keys)
			for _, key := range keys {
				cacheMapVal := cache.cacheMap[key]
				valStruct := valStruct{
					key:    key,
					next:   cacheMapVal.next,
					prev:   cacheMapVal.prev,
					values: cacheMapVal.info,
				}
				gotVals = append(gotVals, valStruct)

			}
			AssertReflectDeepEqual(gotVals, test.expectedVals, t)
		})
	}

}

func TestWriteAndClearingCache(t *testing.T) {

	testDuration := 1 * time.Millisecond
	//this should make the cache check every 10 seconds
	cache := NewCache(testDuration)
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}
	wg.Add(2)

	go cache.ActivateCacheClear(doneChan)

	for _, testVal := range testVals {
		time.Sleep(600 * time.Millisecond)
		cache.WriteToCache(testVal.key, testVal.next, testVal.prev, testVal.values)
		wg.Done()
	}
	wg.Wait()
	doneChan <- struct{}{}

	got := []valStruct{}
	for key, cacheMapVal := range cache.cacheMap {
		valStruct := valStruct{
			key:    key,
			next:   cacheMapVal.next,
			prev:   cacheMapVal.prev,
			values: cacheMapVal.info,
		}
		got = append(got, valStruct)
	}
	expected := []valStruct{
		testVals[1],
	}

	AssertReflectDeepEqual(got, expected, t)
}
