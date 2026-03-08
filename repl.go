package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"math/rand"

	"github.com/Joseap1996/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, ...string) error
}
type config struct {
	pokeapiClient pokeapi.Client
	Next          string
	Previous      string
	LocationName  string
	PokemonName   string
	caughtPokemon map[string]pokeapi.Pokemon
}

var commands map[string]cliCommand

func commandExit(cfg *config, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
func commandHelp(cfg *config, args ...string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	for name, cmd := range commands {
		fmt.Printf("%v: %v\n", name, cmd.description)
	}
	return nil
}
func commandMap(cfg *config, args ...string) error {
	data, err := cfg.pokeapiClient.GetLocationsAreas(cfg.Next)
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
func commandMapb(cfg *config, args ...string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	data, err := cfg.pokeapiClient.GetLocationsAreas(cfg.Previous)
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

func commandExplore(cfg *config, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no location name")
	}
	locationName := args[0]
	data, err := cfg.pokeapiClient.GetLocationAreaPokemon(locationName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %v...\n", locationName)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range data.Pokemon_Encounters {
		fmt.Printf(" - %v\n", pokemon.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("no pokemon to catch")
	}
	pokemonName := args[0]
	data, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}
	catch := true
	roll := rand.Intn(data.BaseExperience)
	fmt.Printf("Throwing a Pokeball at %v...\n", data.Name)

	if roll > 40 {
		catch = false
	}

	if catch {
		fmt.Printf("%v was caught!\n", data.Name)
		cfg.caughtPokemon[data.Name] = data
	} else {
		fmt.Printf("%v escaped!\n", data.Name)
	}
	return nil
}

func commandInspect(cfg *config, args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("no pokemon to inspect")
	}
	pokemonName := args[0]
	if pokemon, ok := cfg.caughtPokemon[pokemonName]; !ok {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %v\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeInfo := range pokemon.Types {
			fmt.Printf("  - %v\n", typeInfo.Type.Name)
		}
	}
	return nil
}

func cleanInput(text string) []string {
	input := strings.ToLower(text)
	words := strings.Fields(input)

	return words
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		userInput := scanner.Text()

		words := cleanInput(userInput)
		if len(words) == 0 {
			continue
		}

		userCommand := words[0]
		args := words[1:]
		cmd, exists := commands[userCommand] // this checks that the command the user typed exist in the commands map
		if exists {
			err := cmd.callback(cfg, args...)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
