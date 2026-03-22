package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fds66/pokedexcli/internal/pokecache"
)

type Client struct {
	httpClient http.Client
}

func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}

}

type PokemonResponse struct {
	Count    int               `json:"count"`
	Next     *string           `json:"next"`
	Previous *string           `json:"previous"`
	Results  []PokemonLocation `json:"results"`
}

type PokemonLocation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (c *Client) GetLocationList(url string, cache *pokecache.Cache) (PokemonResponse, error) {
	body, err := c.GetAPIdata(url, cache)
	if err != nil {
		return PokemonResponse{}, err
	}
	//fmt.Println(string(body)) //if I want to check the body while debugging

	var pokeData PokemonResponse
	if err := json.Unmarshal(body, &pokeData); err != nil {
		return PokemonResponse{}, err
	}

	return pokeData, nil
}

type PokemonList struct {
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`

	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetPokemonList(url string, cache *pokecache.Cache) (PokemonList, error) {
	body, err := c.GetAPIdata(url, cache)
	if err != nil {
		return PokemonList{}, err
	}
	//fmt.Println(string(body)) //if I want to check the body while debugging

	var pokeData PokemonList
	if err := json.Unmarshal(body, &pokeData); err != nil {
		return PokemonList{}, err
	}

	return pokeData, nil
}

func (c *Client) GetAPIdata(url string, cache *pokecache.Cache) ([]byte, error) {

	// Check if the results already exist in the cache

	var body []byte
	cachedData, success := cache.Get(url)
	if success {
		body = cachedData
		fmt.Println("found cached page")
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return body, err
		}

		results, err := c.httpClient.Do(req)
		if err != nil {
			return body, err
		}
		defer results.Body.Close()
		body, err = io.ReadAll(results.Body)
		if err != nil {
			return body, err
		}

		if !cache.Add(url, body) {
			fmt.Println("data failed to add to cache")
		}
	}

	//fmt.Println(string(body)) //if I want to check the body while debugging

	return body, nil
}
