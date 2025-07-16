package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"net/http"
	"encoding/json"
	"io"
	"errors"
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
		"explore": {
			name:		 "explore",
			description: "Shows a list of all the Pokemon in the provided map location",
			callback:	 commandExplore,
		},
	}
}


type cliCommand struct {
	name		string
	description string
	callback 	func(c *Config, cache *pokecache.Cache, s string) (*Config, error)
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

type LocationData struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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
		parameter := ""
		if len(wordSlice) > 1 {
			parameter = wordSlice[1]
		}
		word, ok := getCommands()[wordSlice[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		new_c, err := word.callback(c, cache, parameter)
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

func commandExit(c *Config, _ *pokecache.Cache, s string) (*Config, error) {
	if len(s) > 0 {
		fmt.Println("Invalid command - Exit does not accept additional parameters.")
		return c, nil
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return c, nil
}

func commandHelp(c *Config, _ *pokecache.Cache, s string) (*Config, error) {
	if len(s) > 0 {
		fmt.Println("Help command does not accept additional parameters - displaying help menu.")
	}
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for key, value := range getCommands() {
		fmt.Printf("%s: %s\n", key, value.description)
	}
	return c, nil
}

func commandMap(c *Config, cache *pokecache.Cache, s string) (*Config, error) {
	if len(s) > 0 {
		fmt.Println("Invalid command - Map does not accept additional parameters.")
		return c, nil
	}
	c, err := getLocations(c, cache)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c, err
	}
	for _, result := range c.Results {
		fmt.Printf("%s\n", result.Name)
	}
	return c, nil
}

func commandMapb(c *Config, cache *pokecache.Cache, s string) (*Config, error) {
	if len(s) > 0 {
		fmt.Println("Invalid command - Map does not accept additional parameters.")
		return c, nil
	}
	if c.Previous == nil {
		fmt.Println("You're on the first page.")
		return c, nil
	}
	c, err := getLocationsb(c, cache)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return c, err
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
	}
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
	}
	err := json.Unmarshal(body, &new_c)
	if err != nil {
		return &new_c, err
	}
	return &new_c, nil
}

func commandExplore(c *Config, cache *pokecache.Cache, s string) (*Config, error) {
	if len(s) == 0 {
		err := errors.New("You must provide a location parameter.")
		return c, err
	}
	loc, err := getLocData(c, cache, s)
	if err != nil {
		return c, err
	}
	for _, result := range loc.PokemonEncounters {
		fmt.Printf("%s\n", result.Pokemon.Name)
	}
	return c, err
}

func getLocData(_ *Config, cache *pokecache.Cache, s string) (*LocationData, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + s + "/"
	var loc LocationData
	var body []byte
	val, ok := cache.Get(url)
	if ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return &loc, err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return &loc, err
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			return &loc, fmt.Errorf("Error: Status Code %v", res.StatusCode)
		}
		cache.Add(url, body)
	}
	err := json.Unmarshal(body, &loc)
	if err != nil {
		return &loc, err
	}
	return &loc, nil
}