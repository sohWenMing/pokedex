package config

import (
	"io"
	"net/http"
	"time"

	caughtpokemon "github.com/sohWenMing/pokedex_cli/caught_pokemon"
	"github.com/sohWenMing/pokedex_cli/internal"
)

type Config struct {
	loc_area_offset int
	client          *http.Client
	Writer          io.Writer
	cache           *internal.Cache
	caughtPokemon   *caughtpokemon.CaughtPokemon
}

func InitConfig(w io.Writer, tickerInterval, cacheValidity time.Duration) (config Config, err error) {
	cache, err := internal.InitCache(tickerInterval, cacheValidity)
	if err != nil {
		return Config{}, err
	}
	return Config{-20, nil, w, cache, caughtpokemon.InitCaughtPokemon()}, nil
}

func (c *Config) IncOffset() {
	c.loc_area_offset += 20
	return
}
func (c *Config) DecOffSet() {
	c.loc_area_offset -= 20
	return
}

func (c *Config) ResetOffSet() {
	c.loc_area_offset = 0
}
func (c *Config) GetOffSet() int {
	return c.loc_area_offset
}
func (c *Config) SetClient(initClient *http.Client) {
	c.client = initClient
	return
}
func (c *Config) GetClient() *http.Client {
	return c.client
}

func (c *Config) GetCache() *internal.Cache {
	return c.cache
}

func (c *Config) GetCaughtPokemon() *caughtpokemon.CaughtPokemon {
	return c.caughtPokemon
}
