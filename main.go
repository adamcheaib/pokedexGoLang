package main

import (
	"bufio"
	"fmt"
	"github.com/adamcheaib/pokedexGoLang/pokecache"
	"os"
	"time"
)

var Cache *pokecache.Cache = pokecache.NewCache(time.Duration(3) * time.Second)

type PokemonApiData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

var Commands = map[string]CliCommands{}

func main() {
	navigation := Config{
		next:     "",
		previous: "",
	}

	scanner := bufio.NewScanner(os.Stdin)

	Commands["exit"] =
		CliCommands{
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		}

	Commands["help"] = CliCommands{
		name:        "help",
		description: "Displays a help message",
		callback:    CommandHelp,
	}

	Commands["map"] = CliCommands{
		name:        "map",
		description: "Pokemon map",
		callback:    CommandMap,
	}

	Commands["mapb"] = CliCommands{
		name:        "MapB",
		description: "Go back to the previous location",
		callback:    MapB,
	}

	for {
		fmt.Print("Pokedéx > ")
		scanner.Scan()
		input := CleanInput(scanner.Text())
		requestedCommand, exists := Commands[input[0]]

		if !exists {
			fmt.Println("Command not found")
		} else {
			requestedCommand.callback(&navigation)
		}
	}

}
