package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Data struct {
	Elevation float64        `json:"elevation"`
	Hourly    map[string]any `json:"hourly"`
}

func main() {
	endpoint := "https://api.open-meteo.com/v1/forecast?latitude=52.52&longitude=13.41&hourly=temperature_2m"

	// req, err := http.NewRequest("GET", endpoint, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	client := &http.Client{}
	resp, err := client.Get(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	var data Data
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)
}
