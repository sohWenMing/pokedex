package internal

import (
	"reflect"
	"testing"
)

var cache *Cache

func TestMain(m *testing.M) {
	cache = InitCache()
	m.Run()

}
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
	}

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
