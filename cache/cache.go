package cache

import (
	"sync"
	"time"
)

/*
first build out a structure that maps a url to a list of strings
*/

type cacheMapVal struct {
	info     []string
	cachedOn time.Time
}
type cacheMap map[string]cacheMapVal

type Cache struct {
	cacheValidity time.Duration
	cacheMap      cacheMap
	mu            sync.Mutex
}

func NewCache(cacheValidity time.Duration) *Cache {
	cacheMap := make(map[string]cacheMapVal)
	returnedCache := Cache{
		cacheValidity: cacheValidity,
		cacheMap:      cacheMap,
		mu:            sync.Mutex{},
	}
	return &returnedCache
}

func ActivateCacheClear(c *Cache, tickerChan, counterChan, doneChan chan struct{}) {

	for {
		select {
		case <-doneChan:
			close(tickerChan)
			close(counterChan)
			return
		case <-tickerChan:
			c.mu.Lock()
			for key, val := range c.cacheMap {
				if time.Since(val.cachedOn) > c.cacheValidity {
					delete(c.cacheMap, key)
				}
			}
			counterChan <- struct{}{}
			c.mu.Unlock()
		}
	}
}

func WriteToCache(c *Cache, url string, values []string) {
	cacheMapVal := cacheMapVal{
		info:     values,
		cachedOn: time.Now(),
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[url] = cacheMapVal

}
