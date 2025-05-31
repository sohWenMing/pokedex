package internal

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

var cache *Cache

func TestWriteToCache(t *testing.T) {
	type test struct {
		values      []string
		expected    []string
		key         string
		name        string
		errExpected bool
	}

	tests := []test{
		{
			name: "basic write test",
			key:  "key1",
			values: []string{
				"one", "two", "three",
			},
			expected: []string{
				"one", "two", "three",
			},
			errExpected: false,
		},
		{
			name: "over write test",
			key:  "key1",
			values: []string{
				"two", "three", "four",
			},
			expected: []string{
				"two", "three", "four",
			},
			errExpected: false,
		},
	}
	cache, err := InitCache(5*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	defer cache.AccessTicker().Stop()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := cache.WriteToCache(test.key, test.values)
			if err != nil {
				t.Errorf("didn't expect error, got %v", err)
				return
			}
			cacheEntry := cache.cacheMap[test.key]
			got := cacheEntry.WriteBufToStrings()
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("\ngot: %v\nwant: %v", got, test.expected)
			}
		})
	}
}

func TestClearCache(t *testing.T) {
	cache, err := InitCache(2*time.Second, 1*time.Second)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}

	go func(cache *Cache) {
		for i := range 5 {
			key := fmt.Sprintf("key%d", i)
			vals := make([]string, 0, 5)
			for j := range i + 5 {
				val := fmt.Sprintf("val%d", j)
				vals = append(vals, val)
			}
			cache.WriteToCache(key, vals)
			time.Sleep(500 * time.Millisecond)
		}
	}(cache)

	time.Sleep(2 * time.Second)
	numRecordsAfterSleep := len(cache.cacheMap)

	stopperChannel := make(chan bool)
	go func(stopperChannel chan<- bool) {
		time.Sleep(3 * time.Second)
		stopperChannel <- true
	}(stopperChannel)

	for {
		select {
		case <-cache.clearCacheTicker.C:
			cache.RemoveOutdated()
		case <-stopperChannel:
			cache.AccessTicker().Stop()
			numRecordsAfterClearing := len(cache.cacheMap)
			if numRecordsAfterSleep >= numRecordsAfterClearing {
				t.Errorf("After sleep should be less than after clearing")
				fmt.Println("after sleep: ", numRecordsAfterSleep)
				fmt.Println("after clearing: ", numRecordsAfterClearing)
			}
			return
		}
	}

}
