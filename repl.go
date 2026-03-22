package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/fds66/pokedexcli/internal/pokeapi"
	"github.com/fds66/pokedexcli/internal/pokecache"
)

type Config struct {
	PokeapiClient pokeapi.Client
	Cache         *pokecache.Cache
	Next          *string
	Previous      *string
	Pokedex       map[string]Pokemon
}

func startRepl(configuration *Config) {
	/* start up cli*/
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex >")
		err := scanner.Scan()
		if err == false {
			fmt.Println(scanner.Err())
			break
		}

		inputString := scanner.Text()
		// clean up the input and use the first word as the command
		cleanInputList := cleanInput(inputString)
		if len(cleanInputList) == 0 {
			continue
		}
		commandWord := cleanInputList[0]
		var argument string
		if len(cleanInputList) > 1 {
			argument = cleanInputList[1]

		} else {
			argument = ""
		}
		//if the command exists then execute the callback function

		cmd, exists := getCommands()[commandWord]
		if !exists {
			fmt.Println("Unknown Command")
			continue
		}

		err2 := cmd.callback(configuration, argument)
		if err2 != nil {
			fmt.Printf("Error executing command, %v\n", err)
		}

	}
}

func cleanInput(text string) []string {
	result := []string{}
	result = strings.Fields(text)
	for i, word := range result {
		result[i] = strings.ToLower(word)
	}
	return result
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{

		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Prints next list of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Prints the previous list of locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Prints the pokemon found in a specified location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a specified pokemon",
			callback:    commandCatch,
		},
	}

}

func commandExit(configuration *Config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configuration *Config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	var outputString string

	for _, cmd := range getCommands() {
		outputString = fmt.Sprintf("%s:\t%s", cmd.name, cmd.description)
		fmt.Println(outputString)
	}
	return nil

}

func commandMap(configuration *Config, args ...string) error {
	// displays the names of 20 location areas in Pokemon world. Each new call gives next 20
	baseurl := "https://pokeapi.co/api/v2/location-area/"
	nextPage := configuration.Next
	url := baseurl
	if nextPage != nil {
		url = *nextPage
	}

	pokeData, err := configuration.PokeapiClient.GetLocationList(url, configuration.Cache)
	if err != nil {
		fmt.Printf("API call failed %v", err)
	}
	usePokeData(pokeData, configuration)

	return nil

}

func commandMapb(configuration *Config, args ...string) error {
	// displays the names of 20 location areas in Pokemon world. This call gives the previous 20 if they exist otherwise just tells you are on the first page
	previousPage := configuration.Previous
	if previousPage == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url := *previousPage
	pokeData, err := configuration.PokeapiClient.GetLocationList(url, configuration.Cache)
	if err != nil {
		fmt.Printf("API call failed %v", err)
	}
	usePokeData(pokeData, configuration)
	return nil
}

func usePokeData(pokeData pokeapi.PokemonResponse, configuration *Config) error {
	//location := pokeData.Results[0].Name
	for _, location := range pokeData.Results {
		fmt.Println(location.Name)
	}

	if pokeData.Next != nil {
		configuration.Next = pokeData.Next
	} else {
		configuration.Next = nil
	}
	if pokeData.Previous != nil {
		configuration.Previous = pokeData.Previous
	} else {
		configuration.Previous = nil
	}
	return nil
}

func commandExplore(configuration *Config, args ...string) error {
	// displays the names pokemon in the specified location areas in Pokemon world

	location := args[0]
	fmt.Printf("location %s, args %v\n", location, args)
	if location == "" {
		fmt.Println("no location specified")
		return fmt.Errorf("no location")
	}
	baseurl := "https://pokeapi.co/api/v2/location-area/"
	url := fmt.Sprintf("%s%s", baseurl, location)

	pokeData, err := configuration.PokeapiClient.GetPokemonList(url, configuration.Cache)
	if err != nil {
		fmt.Printf("API call failed %v\n", err)
		return err
	}
	fmt.Printf("Exploring %s...\n", pokeData.Location.Name)
	for _, encounter := range pokeData.PokemonEncounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil

}

type Pokemon struct {
	Name string
	Info pokeapi.PokemonData
}

func commandCatch(configuration *Config, args ...string) error {
	// displays the names pokemon in the specified location areas in Pokemon world

	pokemon := args[0]
	if pokemon == "" {
		fmt.Println("no pokemon specified")
		return fmt.Errorf("no location")
	}
	baseurl := "https://pokeapi.co/api/v2/pokemon/"
	url := fmt.Sprintf("%s%s", baseurl, pokemon)

	pokeData, err := configuration.PokeapiClient.GetPokemonData(url, configuration.Cache)
	if err != nil {
		fmt.Printf("API call failed %v", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokeData.Name)
	baseExperience := pokeData.BaseExperience
	maxExperience := 300

	chance := rand.Float32()

	fmt.Printf("chance %.2f, base experience %v\n", chance, baseExperience)
	if chance*float32(maxExperience) > float32(baseExperience) {
		fmt.Printf("You caught %s, adding to your pokedex\n", pokeData.Name)
		newPokedexEntry := Pokemon{
			Name: pokeData.Name,
			Info: pokeData,
		}

		configuration.Pokedex[pokeData.Name] = newPokedexEntry
		fmt.Println("Pokemon in the Pokedex")
		for _, name := range configuration.Pokedex {
			fmt.Printf("- %s\n", name.Name)
		}
	} else {
		fmt.Printf("%s broke free and was not captured\n", pokeData.Name)
	}
	return nil

}
