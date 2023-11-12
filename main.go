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

type Sender interface {
	Send(*WeatherData) error
}

type SMSSender struct {
	number string
}

func NewSMSSender(number string) *SMSSender {
	return &SMSSender{
		number,
	}
}

func (s *SMSSender) Send(data *WeatherData) error {
	fmt.Println("Sender weather to number: ", s.number)
	return nil
}

type EmailSender struct {
	email string
}

func NewEmailSender(email string) *EmailSender {
	return &EmailSender{
		email: email,
	}
}

func (es *EmailSender) Send(data *WeatherData) error {
	fmt.Println("Sending email to: ", es.email)
	return nil
}

var (
	pollInterval = time.Second * 5
)

type WeatherData struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

type WPoller struct {
	closech chan struct{}
	senders []Sender
}

func NewWPoller(senders ...Sender) *WPoller {
	return &WPoller{
		closech: make(chan struct{}),
		senders: senders,
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
	// handle the data (store it in DB maybe)
	for _, s := range wp.senders {
		if err := s.Send(data); err != nil {
			fmt.Println(err)
		}
	}
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
	smsSender := NewSMSSender("1231150813")
	emailSender := NewEmailSender("wesley@gmail.com")
	wpoller := NewWPoller(smsSender, emailSender)
	go func() {
		wpoller.start()
	}()

	select {}
}
