package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const openWeatherMapURL = "https://api.openweathermap.org/data/2.5/weather"

type weatherResponse struct {
	Name    string `json:"name"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func getWeather(city string, apiKey string) (*weatherResponse, error) {
	client := resty.New()

	resp, err := client.R().
		SetQueryParams(map[string]string{
			"q":     city,
			"appid": apiKey,
		}).
		Get(openWeatherMapURL)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Request failed with status code: %v", resp.StatusCode())
	}

	var weatherData weatherResponse
	if err := json.Unmarshal(resp.Body(), &weatherData); err != nil {
		return nil, err
	}

	return &weatherData, nil
}

func main() {
	godotenv.Load()

	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if apiKey == "" {
		log.Fatal("OpenWeatherMap API key not set. Set it using the OPENWEATHERMAP_API_KEY environment variable.")
	}

	city := "Curitiba" // Pode ser substituído pela cidade desejada.

	weatherData, err := getWeather(city, apiKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Weather in %s: %s. Temperature: %.2f°C\n", weatherData.Name, weatherData.Weather[0].Description, weatherData.Main.Temp)
}
