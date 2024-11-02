package cache

import (
	"reflect"
	"testing"
	"time"

	errorHelpers "github.com/sohWenMing/pokedex/test_error_helpers"
)

var assertVals = errorHelpers.AssertVals

func TestNewCache(t *testing.T) {
	cache := NewCache(50 * time.Millisecond)
	tickerChan := make(chan struct{})
	counterChan := make(chan struct{})
	doneChan := make(chan struct{})
	ticker := time.NewTicker(100 * time.Millisecond)
	counterVal := 0

	go ActivateCacheClear(cache, tickerChan, counterChan, doneChan)

	go func() {
		for range ticker.C {
			tickerChan <- struct{}{}
		}
	}()

	go func() {
		for range counterChan {
			counterVal++
		}
	}()
	time.Sleep(1000 * time.Millisecond)
	go func() {
		doneChan <- struct{}{}
	}()

	assertVals(counterVal, 10, t)

}

func TestWriteToCache(t *testing.T) {
	cache := NewCache(0 * time.Second)
	type testKeyValue struct {
		key    string
		values []string
	}
	type testStruct struct {
		name          string
		testKeyValues []testKeyValue
		expected      [][]string
	}

	expected_vals_set_1 := []string{
		"value 1", "value 2",
	}
	expected_vals_set_2 := []string{
		"value 3", "value 4",
	}
	test_one_key_values := []testKeyValue{
		{
			key:    "url 1",
			values: expected_vals_set_1,
		},
		{
			key:    "url 2",
			values: expected_vals_set_2,
		},
	}
	tests := []testStruct{
		{
			name:          "basic write test",
			testKeyValues: test_one_key_values,
			expected: [][]string{
				expected_vals_set_1, expected_vals_set_2,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := [][]string{}
			for _, keyVal := range test.testKeyValues {
				WriteToCache(cache, keyVal.key, keyVal.values)
			}
			for _, val := range cache.cacheMap {
				got = append(got, val.info)
			}
			if reflect.DeepEqual(got, test.expected) {
				t.Errorf("\ngot: %v\nwant: %v", got, test.expected)
			}
		})
	}

}
