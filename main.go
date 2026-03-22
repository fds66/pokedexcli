package main

import (
	"time"

	"github.com/fds66/pokedexcli/internal/pokeapi"
	"github.com/fds66/pokedexcli/internal/pokecache"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	newCache := pokecache.NewCache(20 * time.Second)
	newPokedex := make(map[string]Pokemon)
	configuration := &Config{
		PokeapiClient: pokeClient,
		Cache:         newCache,
		Next:          nil,
		Previous:      nil,
		Pokedex:       newPokedex,
	}

	startRepl(configuration)

}
