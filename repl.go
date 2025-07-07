package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	"io"
	"github.com/smwalke83/pokedex/internal/pokecache"
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
		"map": {
			name:		 "map",
			description: "Shows the next 20 map locations",
			callback:	 commandMap,
		},
		"mapb": {
			name: 		 "mapb",
			description: "Shows the previous 20 map locations",
			callback:	 commandMapb,
		},
	}
}


type cliCommand struct {
	name		string
	description string
	callback 	func(c *Config, cache *pokecache.Cache) (*Config, error)
}

type Config struct {
	Count		int		`json:"count"`
	Next		string	`json:"next"`
	Previous	*string	`json:"previous"`
	Results		[]struct {
		Name	string	`json:"name"`
		Url		string	`json:"url"`
	} `json:"results"`
}

func startRepl(cache *pokecache.Cache) {
	c := new(Config)
	scan := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		ok := scan.Scan()
		if !ok {
			fmt.Printf("Error: %v\n", scan.Err())
		}
		input := scan.Text()
		wordSlice := cleanInput(input)
		word, ok := getCommands()[wordSlice[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		new_c, err := word.callback(c, cache)
		if err != nil {
			fmt.Println(err)
		}
		*c = *new_c
	}
}

func cleanInput(text string) []string {
	words := strings.Split(strings.ToLower(strings.TrimSpace(text)), " ")
	return words
}

func commandExit(c *Config, _ *pokecache.Cache) (*Config, error) {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return c, nil
}

func commandHelp(c *Config, _ *pokecache.Cache) (*Config, error) {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return c, nil
}

func commandMap(c *Config, cache *pokecache.Cache) (*Config, error) {
	c, err := getLocations(c, cache)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	for _, result := range c.Results {
		fmt.Printf("%s\n", result.Name)
	}
	return c, nil
}

func commandMapb(c *Config, cache *pokecache.Cache) (*Config, error) {
	if c.Previous == nil {
		fmt.Println("You're on the first page.")
		return c, nil
	}
	c, err := getLocationsb(c, cache)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	for _, result := range c.Results {
		fmt.Printf("%s\n", result.Name)
	}
	return c, nil
}

func getLocations(c *Config, cache *pokecache.Cache) (*Config, error) {
	url := c.Next
	var body []byte
	if c.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area"
	}
	var new_c Config
	val, ok := cache.Get(url)
	fmt.Printf("Get: key=[%s] hit=%v bytes=%d\n", url, ok, len(val))
	if ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return &new_c, err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return &new_c, err
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			return &new_c, fmt.Errorf("Error: Status Code %v", res.StatusCode)
		}
		cache.Add(url, body)
		fmt.Printf("Add: key=[%s] bytes=%d\n", url, len(body))
	}
	fmt.Printf("Cache hit: %v, bytes: %v\n", ok, len(body))
	err := json.Unmarshal(body, &new_c)
	if err != nil {
		return &new_c, err
	}
	return &new_c, nil
}

func getLocationsb(c *Config, cache *pokecache.Cache) (*Config, error) {
	var url string
	var new_c Config
	var body []byte
	if c.Previous == nil {
		return c, nil
	} else {
		url = *c.Previous
	}
	val, ok := cache.Get(url)
	fmt.Printf("Get: key=[%s] hit=%v bytes=%d\n", url, ok, len(val))
	if ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return &new_c, err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return &new_c, err
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			return &new_c, fmt.Errorf("Error: Status Code %v", res.StatusCode)
		}
		cache.Add(url, body)
		fmt.Printf("Add: key=[%s] bytes=%d\n", url, len(body))
	}
	fmt.Printf("Cache hit: %v, bytes: %v\n", ok, len(body))
	err := json.Unmarshal(body, &new_c)
	if err != nil {
		return &new_c, err
	}
	return &new_c, nil
}