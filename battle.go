package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Joseap1996/pokedex/internal/pokeapi"
)

func calculateDamage(attack int, defence int) int {
	damage := attack - defence
	return damage
}

func getPokemonStats(p pokeapi.Pokemon) (hp, attack, def, speed int) {
	for _, s := range p.Stats {
		switch s.Stat.Name {
		case "hp":
			hp = s.BaseStat
		case "attack":
			attack = s.BaseStat
		case "defence":
			def = s.BaseStat
		case "speed":
			speed = s.BaseStat
		}
	}
	return hp, attack, def, speed

}

func listUserPokemon(cfg *config) error {
	if len(cfg.caughtPokemon) == 0 {
		fmt.Println("You haven't caught any Pokemon yet!")
		return nil

	}

	for name := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

func commandRandomBattle(cfg *config, args ...string) error {
	location := args[0]
	loc_data, err := cfg.pokeapiClient.GetLocationAreaPokemon(location)
	if err != nil {
		return fmt.Errorf("invalid location")
	}
	fmt.Println("Your Pokemon:")
	err = listUserPokemon(cfg)
	if err != nil {
		return err
	}

	fmt.Print("Pick your pokemon you want to use for battle: ")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	user_Pokemon := strings.ToLower(strings.TrimSpace(scanner.Text()))

	poke_Data, ok := cfg.caughtPokemon[user_Pokemon]
	if !ok {
		fmt.Println("You havent caught that pokemon")
		return nil
	}

	fmt.Println("Searching for wild pokemon...")
	time.Sleep(3 * time.Second) // small delay to simulate loading/searching of pokemon

	index := rand.Intn(len(loc_data.Pokemon_Encounters))
	pokemon := loc_data.Pokemon_Encounters[index]
	fmt.Printf("A wild %v appears! prepare to battle...\n", pokemon.Pokemon.Name)

	time.Sleep(3 * time.Second)

	fmt.Printf("%v! I choose you!\n", user_Pokemon)

	time.Sleep(3 * time.Second)

	wild_Poke_Data, err := cfg.pokeapiClient.GetPokemon(pokemon.Pokemon.Name) // gets wild pokemon data
	if err != nil {
		return err
	}

	wild_Hp, wild_Attack, wild_Def, wild_Speed := getPokemonStats(wild_Poke_Data) // gets wild pokemon stats
	u_Hp, u_Attack, u_Def, u_Speed := getPokemonStats(poke_Data)                  // gets users pokemon stats

	for u_Hp > 0 && wild_Hp > 0 { // battle loop
		if wild_Speed > u_Speed { // speed check to determine who attacks first
			fmt.Printf("%v attacks first!\n", wild_Poke_Data.Name)
			damage := calculateDamage(wild_Attack, u_Def) // call to calcdmg function
			fmt.Printf("%v takes %v of damage!\n", poke_Data.Name, damage)
			u_Hp -= damage // hp subtraction
			time.Sleep(2 * time.Second)

			if u_Hp <= 0 { // hp check to see if pokemon has fainted
				fmt.Printf("%v has fainted!\n", poke_Data.Name)
				break
			}
			// now slower pokemon attacks
			fmt.Printf("Now its %v turn to attack!\n", poke_Data.Name)
			damage = calculateDamage(u_Attack, wild_Def)
			fmt.Printf("%v takes %v of damage!\n", wild_Poke_Data.Name, damage)
			wild_Hp -= damage

			time.Sleep(2 * time.Second)

			if wild_Hp <= 0 {
				fmt.Printf("%v has fainted!\n", wild_Poke_Data.Name)
				fmt.Printf("Would you like to catch the fainted %v? (y/n): ", wild_Poke_Data.Name) // prompt user if they want to catch the deafeated pokemon

				scanner.Scan()
				response := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if response == "y" {
					cfg.caughtPokemon[wild_Poke_Data.Name] = wild_Poke_Data // add wild pokemon to the users
					fmt.Printf("%v was caught!\n", wild_Poke_Data.Name)
					savePokedex(cfg.caughtPokemon) // call the save function to save progress
					break
				} else {
					fmt.Printf("You let %v back into the wild\n", wild_Poke_Data.Name)
					break
				}
			}

		} else {
			fmt.Printf("%v attacks first!\n", poke_Data.Name)
			damage := calculateDamage(u_Attack, wild_Def)
			fmt.Printf("%v takes %v of damage!\n", wild_Poke_Data.Name, damage)
			wild_Hp -= damage

			time.Sleep(2 * time.Second)

			if wild_Hp <= 0 {
				fmt.Printf("%v has fainted!\n", wild_Poke_Data.Name)
				fmt.Printf("Would you like to catch the fainted %v? (y/n): ", wild_Poke_Data.Name)

				scanner.Scan()
				response := strings.ToLower(strings.TrimSpace(scanner.Text()))
				if response == "y" {
					cfg.caughtPokemon[wild_Poke_Data.Name] = wild_Poke_Data
					fmt.Printf("%v was caught!\n", wild_Poke_Data.Name)
					savePokedex(cfg.caughtPokemon)
					break
				} else {
					fmt.Printf("You let %v back into the wild\n", wild_Poke_Data.Name)
					break
				}
			}

			fmt.Printf("Now its %v turn to attack!\n", wild_Poke_Data.Name)
			damage = calculateDamage(wild_Attack, u_Def)
			fmt.Printf("%v takes %v of damage!\n", poke_Data.Name, damage)
			u_Hp -= damage

			time.Sleep(2 * time.Second)

			if u_Hp <= 0 {
				fmt.Printf("%v has fainted!\n", poke_Data.Name)
				break
			}

		}

	}
	return nil

}
