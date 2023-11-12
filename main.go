package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	endpoint = "https://api.open-meteo.com/v1/forecast"
)

var (
	pollInterval = time.Second * 5
)

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type WPoller struct {
}

func NewWPoller() *WPoller {
	return &WPoller{}
}

func (wp *WPoller) start() {
	fmt.Println("Starting the wpoller")

	ticker := time.NewTicker(pollInterval)
	for {
		data, err := getWeatherResults(52.52, 13.41)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(data)
		<-ticker.C
	}
}

func (wp *WPoller) handleData(data *WeatherData) {

}

func main() {
	wp := NewWPoller()
	wp.start()
}

func getWeatherResults(lat, long float64) (*WeatherData, error) {

	uri := fmt.Sprintf("%s?latitude=%.2f&longitude=%.2f&hourly=temperature_2m", endpoint, lat, long)

	fmt.Println("------------------------------------------")
	fmt.Println(uri)
	resp, err := http.Get(uri)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var data WeatherData
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
