package main

import (
	"strings"
	"fmt"
	"os"
)

type cliCommand struct {
	name string
	description string
	callback func(*Config, string) error
}

type Config struct {
	Next string
	Previous string
}

var mover Config

var Commands = map[string]cliCommand{}

func init() {
	Commands["help"] = cliCommand {
		name: "help",
		description: "Displays a help message",
		callback: commandHelp,
	}
	Commands["exit"] = cliCommand {
		name: "exit",
		description: "Exit the Pokedex",
		callback: commandExit,
	}
	Commands["map"] = cliCommand {
		name: "map",
		description: `Displays the Poke Map!
moves the index foward by 20`,
		callback: commandMap,
	}
	Commands["mapb"] = cliCommand {
		name: "mapb",
		description: `Displays the Poke Map!
moves the index backwards by 20`,
		callback: commandMapb,
	}
	Commands["explore"] = cliCommand {
		name: "explore",
		description: `Usage: explore <area-name>
Shows all the pokemon in the area!`,
		callback: commandExplore,
	}
	Commands["catch"] = cliCommand {
		name: "catch",
		description: `Usage: catch <pokemon-name>
Tries to catch the pokemon!
You can see your caught pokemon in the pokedex`,
		callback: commandCatch,
	}
	Commands["inspect"] = cliCommand {
		name: "inspect",
		description: `Usage: inspect <pokemon-name>
Shows the pokemon's stats if you have caught it before!`,
		callback: commandInspect,
	}
	Commands["pokedex"] = cliCommand {
		name: "pokedex",
		description: `Shows all your captured pokemons!`,
		callback: commandPokedex,
	}
	mover.Previous = ""
	mover.Next = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
}

func cleanInput(text string) []string {
	loweredString := strings.ToLower(text)
	formattedSlice := strings.SplitN(loweredString, " ", -1)
	var returnSlice []string
	for i, str := range formattedSlice {
		if formattedSlice[i] == "" {
			continue
		}
		newstr := strings.Trim(str, " ")
		returnSlice = append(returnSlice, newstr)
	}
	return returnSlice
}

func commandExit(config *Config, nope string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!\n")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, nope string) error {
	fmt.Printf(`
Welcome to the Pokedex!
Usage:

`)
	for command, commandBody := range Commands {
		fmt.Printf("%s: %s\n\n", command, commandBody.description)
	}
	return nil
}

func commandMap(config *Config, nope string) error {
	if config.Next == "" {
		return fmt.Errorf("There is no more next pages")
	}
	strArr, err := mapApiCall(config, true)
	if err != nil {
		return err
	}
	for _, area := range strArr {
		fmt.Printf(area + "\n")
	}
	return nil
}

func commandMapb(config *Config, nope string) error {
	if config.Previous == "" {
		return fmt.Errorf("There is no more previous pages")
	}
	strArr, err := mapApiCall(config, false)
	if err != nil {
		return err
	}
	for _, area := range strArr {
		fmt.Printf(area + "\n")
	}
	return nil
}

func commandExplore(config *Config, area string) error {
	fmt.Printf("Exploring %s...\n", area)
	strArr, err := areaApiCall(area)
	if err != nil {
		return err
	}
	fmt.Printf("Found Pokemon:\n")
	for _, pokemon := range strArr{
		fmt.Printf("- %s\n", pokemon)
	}
	return nil
}

var PokeDex = map[string]PokemonInfo{}

func commandCatch(config *Config, pokemonName string) error {
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	poke, caught, err := pokemonApiCall(pokemonName)
	if err != nil {
		return err
	}
	if !caught {
		fmt.Printf("%s escaped\n", pokemonName)
		return nil
	}
	fmt.Printf("%s was caught\n", pokemonName)

	if _, ok := PokeDex[poke.Name]; !ok {
		PokeDex[poke.Name] = poke
	}
	return nil
}

func commandInspect(config *Config, pokemonName string) error {
	pokeinfo, ok := PokeDex[pokemonName]
	if !ok {
		fmt.Printf("You have not caught that pokemon\n")
		return nil
	}
	fmt.Printf(`Name: %s
Height: %d
Weight: %d
Stats:
  -hp: %d
  -attack: %d
  -defense: %d
  -special-attack: %d
  -special-defense: %d
  -speed: %d
`, pokeinfo.Name, pokeinfo.Height, pokeinfo.Weight, pokeinfo.Stats[0].BaseStat, pokeinfo.Stats[1].BaseStat, pokeinfo.Stats[2].BaseStat, pokeinfo.Stats[3].BaseStat, pokeinfo.Stats[4].BaseStat, pokeinfo.Stats[4].BaseStat )
	fmt.Printf("Types:\n")
	for i := range pokeinfo.Types {
		fmt.Printf("  - %s\n", pokeinfo.Types[i].Type.Name)
	}
	return nil
}

func commandPokedex(config *Config, nope string) error {
	if len(PokeDex) == 0 {
		fmt.Printf("It seems like your Pokedex is empty... Go catch some of the pokemons!\n")
		return nil
	}
	fmt.Printf("Your Pokedex:\n")
	for k := range PokeDex {
		fmt.Printf("  - %s\n", k)
	}
	return nil
}

func callFunc(command []string) error {
	commandStruct, ok:= Commands[command[0]]
	if ok == false {
		fmt.Printf("Unkwown command\n")
		return fmt.Errorf("Unkwown command entered")
	}
	callbackfunc := commandStruct.callback
	if len(command) == 2  {
		if err := callbackfunc(&mover, command[1]); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}
	} else if len(command) == 1 {
		if err := callbackfunc(&mover, "nope"); err != nil {
			fmt.Printf("%v\n", err)
			return err
		}	
	} else {
		fmt.Printf("Unknown Command\n")
	}
	return nil
}