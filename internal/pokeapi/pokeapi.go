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

func (c *Client) GetAPIdata(url string, cache *pokecache.Cache) (PokemonResponse, error) {

	// Check if the results already exist in the cache

	var body []byte
	cachedData, success := cache.Get(url)
	if success {
		body = cachedData
		fmt.Println("found cached page")
	} else {
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return PokemonResponse{}, err
		}

		results, err := c.httpClient.Do(req)
		if err != nil {
			return PokemonResponse{}, err
		}
		defer results.Body.Close()
		body, err = io.ReadAll(results.Body)
		if err != nil {
			return PokemonResponse{}, err
		}
		if cache.Add(url, body) {
			fmt.Println("data added to cache")
		}
	}

	//fmt.Println(string(body)) //if I want to check the body while debugging

	var pokeData PokemonResponse
	if err := json.Unmarshal(body, &pokeData); err != nil {
		return PokemonResponse{}, err
	}

	return pokeData, nil
}
