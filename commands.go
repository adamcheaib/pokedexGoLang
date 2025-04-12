package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type HTTPErrorCodes int

const (
	BadRequest = iota + 400
	Unauthorized
	PaymentRequired
	Forbidden
	NotFound
)

type CliCommands struct {
	name        string
	description string
	callback    func(*Config, string) error
}

type Config struct {
	next     string
	previous string
}

type PokemonLocationData struct {
	LocationName string `json:"name"`
	Pokemons     []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func httpErrorHandler(code int) error {
	switch code {
	case BadRequest:
		return errors.New("bad request")
	case Unauthorized:
		return errors.New("unauthorized")
	case PaymentRequired:
		return errors.New("payment required")
	case Forbidden:
		return errors.New("unauthorized")
	case NotFound:
		return errors.New("not found")
	}

	return nil
}

func CommandExit(configuration *Config, misc string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(configuration *Config, misc string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for command := range Commands {
		fmt.Println(command + ": " + Commands[command].description)
	}
	return nil
}

func CommandMap(configuration *Config, misc string) error {
	client := &http.Client{}
	var myData PokemonApiData
	url := "https://pokeapi.co/api/v2/location-area/"

	if cachedData, found := Cache.Get(url); found {
		// Load the information from here instead!
		jsonDecoder := json.NewDecoder(bytes.NewReader(cachedData))
		jsonDecoder.Decode(&myData)
		configuration.next = myData.Next
		configuration.previous = myData.Previous

		for _, item := range myData.Results {
			fmt.Println(item.Name)
		}
		return nil
	}

	// Create the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Stores and then reads the response body into memory as a []byte
	rawData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Deode the bytes data into JSON data.
	jsonDecoder := json.NewDecoder(bytes.NewReader(rawData))
	if err := jsonDecoder.Decode(&myData); err != nil {
		return err
	}

	// Store the *raw data* into the Cache.
	Cache.Add(url, rawData)

	for _, item := range myData.Results {
		fmt.Println(item.Name)
		Cache.Add(item.Name, []byte(item.Url))
	}

	// Ensure the body is closed when the function stops executing
	defer resp.Body.Close()

	// Extract the next URL from the retreieved cachedData
	if myData.Next != "" {
		configuration.next = myData.Next
	} else {
		configuration.next = ""
	}

	// Extract the previous URL from the retrieved cachedData
	if myData.Previous != "" {
		configuration.previous = myData.Previous
	} else {
		configuration.previous = ""
	}

	return nil
}

func MapB(configuration *Config, input string) error {
	if configuration.previous == "" {
		fmt.Println("You're on the first page!")
		return nil
	}

	var myData PokemonApiData

	data, _ := Cache.Get(configuration.previous)

	jsonDecoder := json.NewDecoder(bytes.NewReader(data))
	if err := jsonDecoder.Decode(&myData); err != nil {
		return err
	}

	configuration.next = myData.Next
	configuration.previous = myData.Previous

	for _, item := range myData.Results {
		fmt.Println(item.Name)
	}

	return nil
}

func ExploreCommand(configuration *Config, input string) error {
	var pokemonLocation PokemonLocationData

	formattedUrl := "https://pokeapi.co/api/v2/location-area/" + input

	locationData, found := Cache.Get(input)
	if !found {
		client := &http.Client{}
		request, err := http.NewRequest("GET", formattedUrl, nil)
		res, err := client.Do(request)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}

		jsonDecoder := json.NewDecoder(bytes.NewReader(body))
		if err := jsonDecoder.Decode(&pokemonLocation); err != nil {
			return err
		}

		Cache.Add(input, body)
		for _, item := range pokemonLocation.Pokemons {
			fmt.Println("-", item.Pokemon.Name)
		}
	} else {

		if err := json.Unmarshal(locationData, &pokemonLocation); err != nil {
			return err
		}

		for _, item := range pokemonLocation.Pokemons {
			fmt.Println("-", item.Pokemon.Name)
		}
	}

	return nil
}

func CatchCommand(configuration *Config, input string) error {
	url := "https://pokeapi.co/api/v2/pokemon/" + strings.TrimSpace(input)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := httpErrorHandler(res.StatusCode); err != nil {
		fmt.Println(err)
		return errors.New("The request did not go through!")
	}

	var aPokemon Pokemon

	fmt.Println(fmt.Sprintf("Throwing a Pokeball at %v...", input))
	jsonDecoder := json.NewDecoder(res.Body)
	if err := jsonDecoder.Decode(&aPokemon); err != nil {
		return err
	}

	breach := float64(aPokemon.BaseExperience) * 0.5
	userChance := rand.Intn(aPokemon.BaseExperience)

	time.Sleep(500 * time.Millisecond)
	if userChance > int(breach) {
		fmt.Println(fmt.Sprintf("%v has now been caught", input))
		fmt.Println("You may now inspect it with the 'inspect' command")

		for i := range aPokemon.Stats {
			// Makes it so that I do not have to venture into "StatData" field to access the name!
			// Check the struct for the JSON-data to understand.
			aPokemon.Stats[i].StatName = aPokemon.Stats[i].StatData.Name
		}

		CaughtPokemons[input] = aPokemon
	} else {
		fmt.Println("You have failed to catch the pokemon", input)
	}

	return nil
}

func InspectCommand(configuration *Config, input string) error {

	pokemon, found := CaughtPokemons[input]
	if !found {
		fmt.Println(fmt.Sprintf("You have not caught %v!", input))
		return errors.New("You have not caught that Pokemon!")
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")

	for i := range pokemon.Stats {
		fmt.Println(fmt.Sprintf("-%v: %v", pokemon.Stats[i].StatName, pokemon.Stats[i].BaseStat))
	}

	fmt.Println("Types: ")

	for _, pokeType := range pokemon.Types {
		fmt.Println("\t- ", pokeType.PokeType.Name)
	}

	return nil
}

func PokedexCommand(configuration *Config, input string) error {

	if len(CaughtPokemons) == 0 {
		fmt.Println("No pokémons have been caught")
		return errors.New("No pokémons have been caught")
	}

	fmt.Println("Your Pokedex:")

	for _, pokemon := range CaughtPokemons {
		fmt.Println("\t-", pokemon.Name)
	}

	return nil
}
