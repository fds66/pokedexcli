package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

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

func GetAPIdata(url string) (PokemonResponse, error) {
	results, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer results.Body.Close()
	body, err := io.ReadAll(results.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body)) if I want to check the body while debugging

	var pokeData PokemonResponse
	if err := json.Unmarshal(body, &pokeData); err != nil {
		log.Fatal(err)
	}

	return pokeData, nil
}
