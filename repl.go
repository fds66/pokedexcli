package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func startRepl() {
	/* start up cli*/
	scanner := bufio.NewScanner(os.Stdin)
	configuration := initialiseConfig()

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

func initialiseConfig() *config {
	config := config{Next: nil, Previous: nil}
	return &config

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

type pokemonResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []pokemonLocation `json:"results"`
}

type pokemonLocation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type config struct {
	Next     *string
	Previous *string
}

func commandMap(configuration *config) error {
	// displays the names of 20 location areas in Pokemon world. Each new call gives next 20
	baseurl := "https://pokeapi.co/api/v2/location-area/"
	nextPage := configuration.Next
	if nextPage == nil {
		getAPIdata(baseurl, configuration)
	} else {
		url := *nextPage
		getAPIdata(url, configuration)
	}

	return nil

}

func commandMapb(configuration *config) error {
	// displays the names of 20 location areas in Pokemon world. Each new call gives next 20
	previousPage := configuration.Previous
	if previousPage == nil {
		fmt.Println("you're on the first page")
	} else {
		url := *previousPage
		getAPIdata(url, configuration)
	}

	return nil
}

func getAPIdata(request string, configuration *config) error {
	results, err := http.Get(request)
	if err != nil {
		log.Fatal(err)
	}
	defer results.Body.Close()
	body, err := io.ReadAll(results.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body)) if I want to check the body while debugging

	var pokeData pokemonResponse
	if err := json.Unmarshal(body, &pokeData); err != nil {
		log.Fatal(err)
	}
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
