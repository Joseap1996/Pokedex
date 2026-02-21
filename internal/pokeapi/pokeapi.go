package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type Data struct {
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

func GetLocationsAreas(url string) (Data, error) {
	if url == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	}

	res, err := http.Get(url)
	if err != nil {
		return Data{}, fmt.Errorf("error creating request: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		return Data{}, fmt.Errorf("Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
	}
	if err != nil {
		return Data{}, err
	}
	data := Data{}
	if err := json.Unmarshal(body, &data); err != nil {
		return Data{}, err
	}
	return data, nil

}
