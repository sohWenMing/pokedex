package internal

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/sohWenMing/pokedex_cli/utils"
)

/*
1. the cache should be of type map[string]cacheEntry

the cacheEntry struct should have:
* data []byte
* createdAt time.Time

if i am going to write to the cache, first i need to check if the actual key exists


*/

type CacheEntry struct {
	createdAt time.Time
	data      *bytes.Buffer
}

func (c *CacheEntry) WriteBufToStrings() []string {
	returned := []string{}
	splitStrings := strings.Split(c.data.String(), "\n")
	for _, val := range splitStrings {
		if val != "" {
			returned = append(returned, val)
		}
	}
	return returned
}

func CreateCacheEntry(values []string) (entry CacheEntry, err error) {
	buf := &bytes.Buffer{}
	for _, value := range values {
		prepString := utils.CleanLineAndAddNewLine(value)
		_, err := buf.WriteString(prepString)
		if err != nil {
			return CacheEntry{}, err
		}
	}
	entry.createdAt = time.Now()

	entry.data = buf
	return entry, nil
}

type Cache struct {
	cacheMap map[string]CacheEntry
	mu       sync.RWMutex
}

func (c *Cache) WriteToCache(key string, values []string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.cacheMap == nil {
		return errors.New("cannot write to a nil map")
	}
	entry, err := CreateCacheEntry(values)
	if err != nil {
		return err
	}
	c.cacheMap[key] = entry
	return nil
}

func InitCache() *Cache {
	cacheMap := make(map[string]CacheEntry)
	cache := &Cache{
		cacheMap: cacheMap,
	}
	return cache
}
