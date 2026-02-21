package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Joseap1996/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}
type config struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for name, cmd := range commands {
		fmt.Printf("%v: %v\n", name, cmd.description)
	}
	return nil
}
func commandMap(cfg *config) error {
	data, err := pokeapi.GetLocationsAreas(cfg.Next)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous

	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
	return nil
}
func commandMapb(cfg *config) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	data, err := pokeapi.GetLocationsAreas(cfg.Previous)
	if err != nil {
		return err
	}
	cfg.Next = data.Next
	cfg.Previous = data.Previous

	for _, loc := range data.Results {
		fmt.Println(loc.Name)
	}
	return nil
}

func main() {
	cfg := &config{}

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
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()
		cleanedInput := cleanInput(userInput)
		userCommand := cleanedInput[0]
		cmd, exists := commands[userCommand] // this checks that the command the user typed exist in the commands map
		if exists {
			cmd.callback(cfg)
		} else {
			fmt.Println("Unknown command")
		}
	}
}
