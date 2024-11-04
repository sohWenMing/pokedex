package cache

import (
	"fmt"
	"sync"
	"time"
)

/*
first build out a structure that maps a url to a list of strings
*/

type cacheMapVal struct {
	next     string
	prev     string
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

func (c *Cache) ActivateCacheClear(doneChan chan struct{}) {

	ticker := time.NewTicker(c.cacheValidity)
	tickerChan := make(chan struct{})
	closeTickerChan := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				tickerChan <- struct{}{}
			case <-doneChan:
				closeTickerChan <- struct{}{}
				return
			}
		}
	}()
	/*
		in each iter of for loop in goroutine:
		- if signal is received from ticker.C, send signal to tickerChan
		- if signal is received from doneChan, sends signal to closeTickerChan
		  and then exit goroutine through return
	*/

	go func() {
		<-closeTickerChan
		close(tickerChan)
		close(closeTickerChan)
		/*
			gorountine waits for the signal from closeTickerchan, and then
			closes the tickerChan channel, and then closes the closeTickerChan
			channel
		*/
	}()
	for {
		select {
		case <-tickerChan:
			c.mu.Lock()
			for key, cacheMapVal := range c.cacheMap {
				if time.Since(cacheMapVal.cachedOn) > c.cacheValidity {
					delete(c.cacheMap, key)
				}
			}
			c.mu.Unlock()

		case <-doneChan:
			return
		}
	}
	/*
		in each iter of for loop:
		- if signal is received from tickerChan, checks validity of cached
		  values and deletes values whose validity has passed
		- if signal is received from doneChan, closes the counterChan
		  and exits the ActivateClearCache function through return
	*/
}

func (c *Cache) WriteToCache(url, next, prev string, values []string) {
	cacheMapVal := cacheMapVal{
		next:     next,
		prev:     prev,
		info:     values,
		cachedOn: time.Now(),
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[url] = cacheMapVal

}

func (c *Cache) GetFromCache(url string) (next, prev string, values []string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cacheMapVal, ok := c.cacheMap[url]
	if !ok {
		return "", "", []string{}, fmt.Errorf("url: %s not found in cache", url)
	}
	return cacheMapVal.next, cacheMapVal.prev, cacheMapVal.info, nil
}
