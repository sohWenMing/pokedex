package caughtpokemon

import (
	"errors"

	structdefinitions "github.com/sohWenMing/pokedex_cli/struct_definitions"
)

type CaughtPokemon struct {
	numCaught int
	pokeMap   map[string]structdefinitions.Pokemon
}

func InitCaughtPokemon() *CaughtPokemon {
	pokeMap := make(map[string]structdefinitions.Pokemon)
	return &CaughtPokemon{
		numCaught: 0,
		pokeMap:   pokeMap,
	}
}

func (c *CaughtPokemon) Add(name string, pokemon structdefinitions.Pokemon) error {
	_, found := c.pokeMap[name]
	if found {
		return errors.New("%s has already been caught!")
	}
	c.pokeMap[name] = pokemon
	c.numCaught += 1
	return nil
}

func (c *CaughtPokemon) Find(name string) (pokemon structdefinitions.Pokemon, isFound bool) {
	pokemon, found := c.pokeMap[name]
	if !found {
		return structdefinitions.Pokemon{}, false
	}
	return pokemon, true
}
func (c *CaughtPokemon) Delete(name string) {
	_, found := c.pokeMap[name]
	if found {
		delete(c.pokeMap, name)
	}
}
