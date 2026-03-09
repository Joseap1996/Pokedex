package main

import (
	"encoding/json"
	"os"

	"github.com/Joseap1996/pokedex/internal/pokeapi"
)

// save function
func savePokedex(pokedex map[string]pokeapi.Pokemon) error {
	data, err := json.Marshal(pokedex)
	if err != nil {
		return err
	}
	return os.WriteFile("pokedex.json", data, 0644)

}

func loadPokedex() (map[string]pokeapi.Pokemon, error) {
	data, err := os.ReadFile("pokedex.json")
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]pokeapi.Pokemon{}, nil
		}
		return nil, err
	}

	var pokedex map[string]pokeapi.Pokemon
	err = json.Unmarshal(data, &pokedex)
	return pokedex, nil

}
