package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Weather struct to hold the response from the Weather API
type Weather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`
	Forecast struct {
		Forecastday []struct {
			Hour []struct {
				TimeEpoch int64   `json:"time_epoch"`
				TempC     float64 `json:"temp_c"`
				Condition struct {
					Text string `json:"text"`
				} `json:"condition"`
				ChanceOfRain float64 `json:"chance_of_rain"`
			} `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

func main() {
	// Default location is Iasi
	q := "Iasi"

	// If a command-line argument is provided, use it as the location
	if len(os.Args) >= 2 {
		q = os.Args[1]
	}

	// Make a GET request to the Weather API with the specified location
	res, err := http.Get("http://api.weatherapi.com/v1/forecast.json?key=ec180872243c4f57a4f153631230105&q=" + q + "&days=1&aqi=no&alerts=no")
	if err != nil {
		panic(err) // Panic if there's an error making the request
	}
	defer res.Body.Close() // Close the response body after we're done with it

	// Check if the response status code is 200 (OK)
	if res.StatusCode != 200 {
		panic("Weather API not available") // Panic if the API is not available
	}

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err) // Panic if there's an error reading the response body
	}

	// Converts the JSON data received from the weather api into the Golang Weather struct
	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err) // Panic if there's an error unmarshaling the JSON data
	}

	// Extract the location, current weather, and hourly forecast from the response
	location, current, hours := weather.Location, weather.Current, weather.Forecast.Forecastday[0].Hour

	// Print the current weather information
	fmt.Printf(
		"%s, %s: %.0fC, %s\n",
		location.Name,
		location.Country,
		current.TempC,
		current.Condition.Text,
	)

	// Print the hourly forecast for the next 24 hours
	for _, hour := range hours {
		date := time.Unix(hour.TimeEpoch, 0) // Convert Unix timestamp to time.Time
		if date.Before(time.Now()) {
			continue // Skip hours in the past
		}

		fmt.Printf(
			"%s - %.0fC, %.0f%%, %s\n",
			date.Format("15:04"), // Format the time as HH:MM
			hour.TempC,
			hour.ChanceOfRain,
			hour.Condition.Text,
		)
	}
}
