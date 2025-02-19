package main

type locations struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []locationArea `json:"results"`
}

type locationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type locationPokemon struct {
	Name               string      `json:"name"`
	PokemonEncounters []encounter `json:"pokemon_encounters"`
}

type encounter struct {
	Pokemon struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"pokemon"`
}
