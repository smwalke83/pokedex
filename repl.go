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
	"math/rand"
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
		"catch": {
			name:		 "catch",
			description: "Throw a pokeball at a pokemon",
			callback:	 commandCatch,
		},
		"inspect": {
			name:		 "inspect",
			description: "Learn about a pokemon in your pokedex",
			callback:	 commandInspect,
		},
		"pokedex": {
			name:		 "pokedex",
			description: "View the pokemon you've added to your pokedex",
			callback:	 commandPokedex,
		},
	}
}


type cliCommand struct {
	name		string
	description string
	callback 	func(c *Config, cache *pokecache.Cache, s string, pokedex map[string]PokeData) (*Config, error)
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

type PokeData struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order int `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			Ability  any  `json:"ability"`
			IsHidden bool `json:"is_hidden"`
			Slot     int  `json:"slot"`
		} `json:"abilities"`
	} `json:"past_abilities"`
}

func startRepl(cache *pokecache.Cache) {
	c := new(Config)
	pokedex := make(map[string]PokeData)
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
		new_c, err := word.callback(c, cache, parameter, pokedex)
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

func commandExit(c *Config, _ *pokecache.Cache, s string, _ map[string]PokeData) (*Config, error) {
	if len(s) > 0 {
		fmt.Println("Invalid command - Exit does not accept additional parameters.")
		return c, nil
	}
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return c, nil
}

func commandHelp(c *Config, _ *pokecache.Cache, s string, _ map[string]PokeData) (*Config, error) {
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

func commandMap(c *Config, cache *pokecache.Cache, s string, _ map[string]PokeData) (*Config, error) {
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

func commandMapb(c *Config, cache *pokecache.Cache, s string, _ map[string]PokeData) (*Config, error) {
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

func commandExplore(c *Config, cache *pokecache.Cache, s string, _ map[string]PokeData) (*Config, error) {
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

func commandCatch(c *Config, cache *pokecache.Cache, s string, pokedex map[string]PokeData) (*Config, error) {
	if len(s) == 0 {
		err := errors.New("Please enter the name of the Pokemon you wish to catch")
		return c, err
	}
	poke, err := getPokeData(cache, s)
	if err != nil {
		return c, err
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", s)
	randomNumber := rand.Intn(poke.BaseExperience)
	if randomNumber < 40 {
		fmt.Printf("%s was caught!\n", s)
		fmt.Printf("You may now inspect it with the inspect command.\n")
		_, ok := pokedex[s]
		if !ok {
			pokedex[s] = *poke
		}
	} else {
		fmt.Printf("%s escaped!\n", s)
	}
	return c, nil
}

func getPokeData(cache *pokecache.Cache, s string) (*PokeData, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + s + "/"
	var poke PokeData
	var body []byte
	val, ok := cache.Get(url)
	if ok {
		body = val
	} else {
		res, err := http.Get(url)
		if err != nil {
			return &poke, err
		}
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return &poke, err
		}
		res.Body.Close()
		if res.StatusCode > 299 {
			return &poke, fmt.Errorf("Error: Status Code %v", res.StatusCode)
		}
		cache.Add(url, body)
	}
	err := json.Unmarshal(body, &poke)
	if err != nil {
		return &poke, err
	}
	return &poke, nil
}

func commandInspect(c *Config, _ *pokecache.Cache, s string, pokedex map[string]PokeData) (*Config, error) {
	pokemon, ok := pokedex[s]
	if !ok {
		fmt.Printf("you have not caught that pokemon\n")
	} else {
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Printf("Stats:\n")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Printf("Types:\n")
		for _, t := range pokemon.Types {
			fmt.Printf("  -%v\n", t.Type.Name)
		}
	}
	return c, nil
}

func commandPokedex(c *Config, _ *pokecache.Cache, _ string, pokedex map[string]PokeData) (*Config, error) {
	fmt.Println("Your Pokedex:")
	if len(pokedex) == 0 {
		fmt.Println("You haven't caught any pokemon!")
	}
	for key, _ := range pokedex {
		fmt.Printf(" - %s\n", key)
	}
	return c, nil
}