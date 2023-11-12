run: build
	./bin/weather-polling

build: 
	go build -o bin/weather-polling

