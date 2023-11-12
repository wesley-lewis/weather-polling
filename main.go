package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m"
)

type Data struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

func main() {

}

func getWeatherResults(lat, long float64) (*Data, error) {

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	fmt.Println(data)
	return &data, nil
}
