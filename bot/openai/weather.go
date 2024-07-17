package openai

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func checkWeather(location string) string {
	fmt.Println("Got location: ", location)
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return "WEATHER_API_KEY is not set"
	}

	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=yes&alerts=yes", apiKey, location)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("Failed to fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Sprintf("Failed to read weather data: %v", err)
	}

	var currentWeatherData struct {
		Location struct {
			Name      string `json:"name"`
			Region    string `json:"region"`
			Country   string `json:"country"`
			LocalTime string `json:"localtime"`
		} `json:"location"`
		Current struct {
			TempC     float64 `json:"temp_c"`
			Condition struct {
				Text string `json:"text"`
			} `json:"condition"`
			WindKph      float64 `json:"wind_kph"`
			Humidity     int     `json:"humidity"`
			FeelsLikeC   float64 `json:"feelslike_c"`
			PressureMb   float64 `json:"pressure_mb"`
			VisibilityKm float64 `json:"vis_km"`
			UV           float64 `json:"uv"`
		} `json:"current"`
		Alerts struct {
			Alert []struct {
				Event       string `json:"event"`
				Description string `json:"desc"`
			} `json:"alert"`
		} `json:"alerts"`
	}
	if err := json.Unmarshal(body, &currentWeatherData); err != nil {
		return fmt.Sprintf("Failed to parse weather data: %v", err)
	}

	weatherMsg := fmt.Sprintf("Weather in %s, %s (%s) | Local Time: %s | %s | Temp: %0.1f°C | Wind: %0.1f kph | Humidity: %d%% | Feels Like: %0.1f°C | Pressure: %0.1f mb | Visibility: %0.1f km | UV Index: %0.1f",
		currentWeatherData.Location.Name, currentWeatherData.Location.Region, currentWeatherData.Location.Country, currentWeatherData.Location.LocalTime,
		currentWeatherData.Current.Condition.Text, currentWeatherData.Current.TempC, currentWeatherData.Current.WindKph, currentWeatherData.Current.Humidity,
		currentWeatherData.Current.FeelsLikeC, currentWeatherData.Current.PressureMb, currentWeatherData.Current.VisibilityKm, currentWeatherData.Current.UV)
	if len(currentWeatherData.Alerts.Alert) > 0 {
		for _, alert := range currentWeatherData.Alerts.Alert {
			weatherMsg += fmt.Sprintf(" | Alert: %s - %s", alert.Event, alert.Description)
		}
	}

	fmt.Println("Returning: ", weatherMsg)

	return weatherMsg
}
