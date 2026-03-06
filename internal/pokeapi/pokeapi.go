package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/Joseap1996/pokedex/internal/pokecache"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}
type Pokemon struct {
	Name            string `json:"name"`
	Base_Experience int    `json:"base_experience"`
}
type Location struct {
	Name               string `json:"name"`
	Url                string `json:"url"`
	Pokemon_Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			Url  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

type Data struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

func (c *Client) GetLocationsAreas(url string) (Data, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}
	dat, ok := c.cache.Get(url)
	if ok {
		locationData := Data{}
		err := json.Unmarshal(dat, &locationData)
		if err != nil {
			return Data{}, err
		}

		return locationData, nil

	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Data{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Data{}, err
	}
	defer res.Body.Close()

	dat, err = io.ReadAll(res.Body)
	if err != nil {
		return Data{}, err
	}

	c.cache.Add(url, dat)

	locationData := Data{}
	err = json.Unmarshal(dat, &locationData)
	if err != nil {
		return Data{}, err
	}

	return locationData, nil

}

func (c *Client) GetLocationAreaPokemon(locationName string) (Location, error) {
	url := "https://pokeapi.co/api/v2/location-area/" + locationName
	dat, ok := c.cache.Get(url)
	if ok {
		locationData := Location{}
		err := json.Unmarshal(dat, &locationData)
		if err != nil {
			return Location{}, err
		}
		return locationData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Location{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Location{}, err
	}
	defer res.Body.Close()

	dat, err = io.ReadAll(res.Body)
	if err != nil {
		return Location{}, err
	}

	c.cache.Add(url, dat)

	locationData := Location{}
	err = json.Unmarshal(dat, &locationData)
	if err != nil {
		return Location{}, err
	}

	return locationData, nil

}
func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	dat, ok := c.cache.Get(url)
	if ok {
		pokemonData := Pokemon{}
		err := json.Unmarshal(dat, &pokemonData)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonData, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}

	defer res.Body.Close()

	dat, err = io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, dat)

	pokemonData := Pokemon{}
	err = json.Unmarshal(dat, &pokemonData)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemonData, nil
}
