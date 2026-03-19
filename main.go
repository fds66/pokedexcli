package main

import (
	"time"

	"github.com/fds66/pokedexcli/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5 * time.Second)
	configuration := &config{
		PokeapiClient: pokeClient,
		Next:          nil,
		Previous:      nil,
	}
	startRepl(configuration)

}
