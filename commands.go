package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type CliCommands struct {
	name        string
	description string
	callback    func(*Config) error
}

type Config struct {
	next     string
	previous string
}

func CommandExit(configuration *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func CommandHelp(configuration *Config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	for command := range Commands {
		fmt.Println(command + ": " + Commands[command].description)
	}
	return nil
}

func CommandMap(configuration *Config) error {
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

func MapB(configuration *Config) error {
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
