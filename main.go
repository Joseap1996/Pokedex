package main

import (
	"time"

	"github.com/Joseap1996/pokedex/internal/pokeapi"
)

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)
	cfg := &config{
		pokeapiClient: pokeClient,
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}

	commands = map[string]cliCommand{
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
			description: "Displays the names of 20 location areas in the Pokemon world",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays list of pokemon in a area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catches a pokemon and adds them to the pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "shows caught pokemon information",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "lists caught pokemon",
			callback:    commandPokedex,
		},
		"delete": {
			name:        "delete",
			description: "deletes current save file",
			callback:    commandDelete,
		},
	}

	startRepl(cfg)
}
