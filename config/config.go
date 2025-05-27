package config

import "net/http"

type Config struct {
	loc_area_offset int
	client          *http.Client
}

func InitConfig() Config {
	return Config{-20, nil}
}

func (c *Config) IncOffset() {
	c.loc_area_offset += 20
	return
}
func (c *Config) DecOffSet() {
	c.loc_area_offset -= 20
	return
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

/*
we want to be able to keep track of which page it's on ...
it will always give me 20 records

for map command, we will increment the config, and then call the API
for the mapb command, we first check if the config is <= 0, if it is
then we just show you're on the first page

else: decrement config by 20, and then call the map

*/
