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

func MapCaughtPokemonToInspect(pokemon structdefinitions.Pokemon) (inspectStruct structdefinitions.PokemonInspectStruct) {
	name := pokemon.Name
	height := pokemon.Height

	hp := 0
	attack := 0
	defense := 0
	specialAttack := 0
	speed := 0

	for _, stat := range pokemon.Stats {
		switch stat.Stat.Name {
		case "hp":
			hp = stat.BaseStat
		case "attack":
			attack = stat.BaseStat
		case "defense":
			defense = stat.BaseStat
		case "special-attack":
			specialAttack = stat.BaseStat
		case "speed":
			speed = stat.BaseStat
		}
	}

	pokemonTypes := make([]string, 0, len(pokemon.Types))
	for _, pokemonType := range pokemon.Types {
		pokemonTypes = append(pokemonTypes, pokemonType.Type.Name)
	}

	return structdefinitions.PokemonInspectStruct{
		Name:   name,
		Height: height,
		Stats: structdefinitions.PokemonStats{
			Hp:            hp,
			Attack:        attack,
			Defense:       defense,
			SpecialAttack: specialAttack,
			Speed:         speed,
		},
		Types: pokemonTypes,
	}
}

type PokemonInspectStruct struct {
	Name   string
	Height int
	Stats  struct {
		Hp            int
		Attack        int
		Defense       int
		SpecialAttack int
		Speed         int
	}
	types []string
}
