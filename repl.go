package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fds66/pokedexcli/internal/pokeapi"
)

type config struct {
	PokeapiClient pokeapi.Client
	Next          *string
	Previous      *string
}

func startRepl(configuration *config) {
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
		//if the command exists then execute the callback function

		cmd, exists := getCommands()[commandWord]
		if !exists {
			fmt.Println("Unknown Command")
			continue
		}

		err2 := cmd.callback(configuration)
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
	callback    func(*config) error
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
	}

}

func commandExit(configuration *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(configuration *config) error {
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

func commandMap(configuration *config) error {
	// displays the names of 20 location areas in Pokemon world. Each new call gives next 20
	baseurl := "https://pokeapi.co/api/v2/location-area/"
	nextPage := configuration.Next
	url := baseurl
	if nextPage != nil {
		url = *nextPage
	}
	pokeData, err := configuration.PokeapiClient.GetAPIdata(url)
	if err != nil {
		fmt.Printf("API call failed %v", err)
	}
	usePokeData(pokeData, configuration)

	return nil

}

func commandMapb(configuration *config) error {
	// displays the names of 20 location areas in Pokemon world. This call gives the previous 20 if they exist otherwise just tells you are on the first page
	previousPage := configuration.Previous
	if previousPage == nil {
		fmt.Println("you're on the first page")
		return nil
	}
	url := *previousPage
	pokeData, err := configuration.PokeapiClient.GetAPIdata(url)
	if err != nil {
		fmt.Printf("API call failed %v", err)
	}
	usePokeData(pokeData, configuration)
	return nil
}

func usePokeData(pokeData pokeapi.PokemonResponse, configuration *config) error {
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
