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
	closech chan struct{}
}

func NewWPoller() *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
	}
}

func (wp *WPoller) close() {
	close(wp.closech)
}

func (wp *WPoller) start() {
	fmt.Println("Starting the wpoller")

	ticker := time.NewTicker(pollInterval)
outer:
	for {
		select {
		case <-ticker.C:
			data, err := getWeatherResults(52.52, 13.41)

			if err != nil {
				log.Fatal(err)
			}

			if err := wp.handleData(data); err != nil {
				log.Fatal(err)
			}
		case <-wp.closech:
			// handle the graceful shutdown
			break outer
		}
	}

	fmt.Println("wpoller stopped gracefully")
}

func (wp *WPoller) handleData(data *WeatherData) error {
	fmt.Println(data)

	return nil
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

func main() {
	wpoller := NewWPoller()
	go func() {
		wpoller.start()
	}()

	time.Sleep(time.Second * 3)

	wpoller.close()
	select {}
}
