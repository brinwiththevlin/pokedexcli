package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/brinwiththevlin/pokedexcli/internal/pokecache"
)

type config struct {
	nextURL     string
	previousURL string
	cache       *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "display the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "list all encounters possible at this location",
			callback:    commandExplore,
		},
	}
}

func commandExplore(c *config, params ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + params[0]

	var data []byte
	if stream, ok := c.cache.Get(url); ok {
		data = stream
	} else {
		var err error
		data, err = httpGet(url)
		if err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	localPokemon := locationPokemon{}
	err := json.Unmarshal(data, &localPokemon)
	if err != nil {
		return err
	}

	fmt.Println("Pokemon in this location:")
	for _, p := range localPokemon.PokemonEncounters {
		fmt.Printf("- %s\n", p.Pokemon.Name)
	}

	return nil
}

func commandMapb(c *config, params ...string) error {
	var url string
	if c.previousURL == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	url = c.previousURL

	var data []byte
	if stream, ok := c.cache.Get(url); ok {
		data = stream
	} else {
		var err error
		data, err = httpGet(url)
		if err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	locs := locations{}
	err := json.Unmarshal(data, &locs)
	if err != nil {
		return err
	}

	c.nextURL = locs.Next
	c.previousURL = locs.Previous

	for _, r := range locs.Results {
		fmt.Println(r.Name)
	}

	return nil
}

func commandMap(c *config, params ...string) error {
	var url string
	if c.nextURL == "" {
		url = "https://pokeapi.co/api/v2/location-area?limit=20"
	} else {
		url = c.nextURL
	}
	var data []byte
	if stream, ok := c.cache.Get(url); ok {
		data = stream
	} else {
		var err error
		data, err = httpGet(url)
		if err != nil {
			return err
		}
		c.cache.Add(url, data)
	}

	locs := locations{}
	err := json.Unmarshal(data, &locs)
	if err != nil {
		return err
	}

	c.nextURL = locs.Next
	c.previousURL = locs.Previous

	for _, r := range locs.Results {
		fmt.Println(r.Name)
	}

	return nil
}
func commandExit(*config, ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(*config, ...string) error {
	output := "Welcome to the Pokedex!\nUsage:\n\n"
	for _, c := range getCommands() { // Call `getCommands` here
		output += fmt.Sprintf("%s: %s\n", c.name, c.description)
	}
	fmt.Print(output + "\n")
	return nil
}

func httpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d\n", res.StatusCode)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
