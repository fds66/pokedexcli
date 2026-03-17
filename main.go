package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		err := scanner.Scan()
		if err == false {
			fmt.Println(scanner.Err())
			break
		}

		inputString := scanner.Text()
		inputString = strings.ToLower(inputString)
		inputStringParts := strings.Fields(inputString)
		if len(inputStringParts) != 0 {
			outputString := fmt.Sprintf("Your command was: %s", inputStringParts[0])
			fmt.Println(outputString)
		}

	}

}
