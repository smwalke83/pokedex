package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
)

func getCommands() map[string]cliCommand {
	return map[string]cliCommand {
		"exit": {
			name:		 "exit",
			description: "Exit the Pokedex",
			callback:	 commandExit,
		},
		"help": {
			name:		 "help",
			description: "Displays a help message",
			callback:	 commandHelp,
		},
	}
}


type cliCommand struct {
	name		string
	description string
	callback 	func() error
}

func startRepl() {
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scan.Scan()
		if !ok {
			fmt.Println("Error: %v", scan.Err())
		}
		input := scan.Text()
		wordSlice := cleanInput(input)
		word, ok := getCommands()[wordSlice[0]]
		if !ok {
			fmt.Println("Unknown command")
		}
		err := word.callback()
		if err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Split(strings.ToLower(strings.TrimSpace(text)), " ")
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return nil
}