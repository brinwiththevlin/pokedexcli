package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/brinwiththevlin/pokedexcli/internal/pokecache"
)

type config struct {
	nextURL     string
	previousURL string
	cache       *pokecache.Cache
	pokedex     map[string]pokemon
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
		"catch": {
			name:        "catch",
			description: "attempt to catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "basic information on a pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "list all caught pokemon",
			callback:    commandPokedex,
		},
	}
}

func fetch(url string, c *config) ([]byte, error) {

	var data []byte
	if stream, ok := c.cache.Get(url); ok {
		data = stream
	} else {
		var err error
		data, err = httpGet(url)
		if err != nil {
			return nil, err
		}
		c.cache.Add(url, data)
	}
	return data, nil
}

func commandPokedex(c *config, params ...string) error {
	fmt.Println("Your Pokedex:")
	for _, p := range c.pokedex {
		fmt.Printf("- %s\n", p.Name)
	}
	return nil
}

func commandInspect(c *config, params ...string) error {
	poke, ok := c.pokedex[params[0]]
	if !ok {
		fmt.Printf("%s was not caught yet!", params[0])
		return nil
	}

	fmt.Printf("Name: %s\n", poke.Name)
	fmt.Printf("Height: %d\n", poke.Height)
	fmt.Printf("Weight: %d\n", poke.Weight)
	fmt.Println("Stats:")
	for _, stat := range poke.Stats {
		fmt.Printf("\t%s: %d\n", stat.Stat.Name, stat.Val)
	}
	fmt.Println("Types:")
	for _, t := range poke.Types {
		fmt.Printf("\t%s\n", t.Type.Name)
	}
	return nil
}

func commandCatch(c *config, params ...string) error {
	if _, ok := c.pokedex[params[0]]; ok {
		fmt.Printf("you have already caught a %s\n", params[0])
		return nil
	}
	url := "https://pokeapi.co/api/v2/pokemon/" + params[0]

	data, err := fetch(url, c)

	if err != nil {
		if err.Error() == errors.New("responce failed with status code: 404").Error() {
			fmt.Printf("%s is not a pokemon\n", params[0])
			return nil
		}
		return err
	}

	poke := pokemon{}
	err = json.Unmarshal(data, &poke)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", poke.Name)
	if rand.Float64() >= 1-(float64(poke.BaseExp)/300) {
		fmt.Printf("%s was caught!\n", poke.Name)
		c.pokedex[poke.Name] = poke
	} else {
		fmt.Printf("%s escaped!\n", poke.Name)
	}

	return nil
}
func commandExplore(c *config, params ...string) error {
	url := "https://pokeapi.co/api/v2/location-area/" + params[0]

	data, err := fetch(url, c)
	if err != nil {
		return err
	}

	localPokemon := locationPokemon{}
	err = json.Unmarshal(data, &localPokemon)
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

	data, err := fetch(url, c)
	if err != nil {
		return err
	}

	locs := locations{}
	err = json.Unmarshal(data, &locs)
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

	data, err := fetch(url, c)
	if err != nil {
		return err
	}

	locs := locations{}
	err = json.Unmarshal(data, &locs)
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
		err := fmt.Errorf("responce failed with status code: %v", res.StatusCode)
		return nil, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
