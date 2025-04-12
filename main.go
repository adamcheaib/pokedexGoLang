package main

import (
	"bufio"
	"fmt"
	"github.com/adamcheaib/pokedexGoLang/pokecache"
	"os"
	"time"
)

var Cache *pokecache.Cache = pokecache.NewCache(time.Duration(30) * time.Second)
var CaughtPokemons = make(map[string]Pokemon)

type PokemonApiData struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"results"`
}

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		StatName string
		StatData struct {
			Name string `json:"name"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		PokeType struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
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

	Commands["explore"] = CliCommands{
		name:        "explore",
		description: "ExploreCommand a specific area and see which Pokémons are there",
		callback:    ExploreCommand,
	}

	Commands["catch"] = CliCommands{
		name:        "catch",
		description: "Attempt to catch a Pokémon",
		callback:    CatchCommand,
	}

	Commands["inspect"] = CliCommands{
		name:        "inspect",
		description: "Inspect the Pokémons that you have caught",
		callback:    InspectCommand,
	}

	Commands["pokedex"] = CliCommands{
		name:        "Pokedex",
		description: "Check the list of all your Pokémons that you have caught",
		callback:    PokedexCommand,
	}

	for {
		fmt.Print("Pokedéx > ")
		scanner.Scan()
		input := CleanInput(scanner.Text())
		var arg string = ""
		if len(input) >= 2 {
			arg = input[1]
		}
		requestedCommand, exists := Commands[input[0]]

		if !exists {
			fmt.Println("Command not found")
		} else {
			requestedCommand.callback(&navigation, arg)
		}
	}

}
