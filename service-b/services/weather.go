package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.opentelemetry.io/otel"
)

type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func GetTemperatureByCity(ctx context.Context, city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return 0, fmt.Errorf("weather API key not set")
	}

	tr := otel.Tracer("service-b")
	ctx, span := tr.Start(ctx, "WeatherAPI Lookup")
	defer span.End()

	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, url.QueryEscape(city))
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("failed to call WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("weatherapi returned status %d: %s", resp.StatusCode, string(body))
	}

	var data WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, fmt.Errorf("failed to decode weather response: %w", err)
	}

	return data.Current.TempC, nil
}
