package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
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

		err2 := cmd.callback()
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
	callback    func() error
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
	}

}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	var outputString string

	for _, cmd := range getCommands() {
		outputString = fmt.Sprintf("%s: %s", cmd.name, cmd.description)
		fmt.Println(outputString)
	}
	return nil

}
