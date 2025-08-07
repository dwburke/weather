package types

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NWS API Response structures
type PointsResponse struct {
	Properties PointsProperties `json:"properties"`
}

type PointsProperties struct {
	GridID           string `json:"gridId"`
	GridX            int    `json:"gridX"`
	GridY            int    `json:"gridY"`
	Forecast         string `json:"forecast"`
	ForecastHourly   string `json:"forecastHourly"`
	ForecastGridData string `json:"forecastGridData"`
}

type ForecastResponse struct {
	Properties ForecastProperties `json:"properties"`
}

type ForecastProperties struct {
	Periods []ForecastPeriod `json:"periods"`
}

type ForecastPeriod struct {
	Number           int    `json:"number"`
	Name             string `json:"name"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	IsDaytime        bool   `json:"isDaytime"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	TemperatureTrend string `json:"temperatureTrend"`
	WindSpeed        string `json:"windSpeed"`
	WindDirection    string `json:"windDirection"`
	Icon             string `json:"icon"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}

// WeatherClient handles NWS API interactions
type WeatherClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		BaseURL: "https://api.weather.gov",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetForecastByCoordinates gets weather forecast for given latitude and longitude
func (w *WeatherClient) GetForecastByCoordinates(lat, lon float64) (*ForecastResponse, error) {
	// First, get the grid information for the coordinates
	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", w.BaseURL, lat, lon)

	req, err := http.NewRequest("GET", pointsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// NWS API requires a User-Agent header
	req.Header.Set("User-Agent", "weather-app/1.0 (your-email@example.com)")

	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get points data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("NWS API error: %s - %s", resp.Status, string(body))
	}

	var pointsResp PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode points response: %w", err)
	}

	// Now get the forecast using the forecast URL from the points response
	forecastReq, err := http.NewRequest("GET", pointsResp.Properties.Forecast, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast request: %w", err)
	}

	forecastReq.Header.Set("User-Agent", "weather-app/1.0 (your-email@example.com)")

	forecastResp, err := w.HTTPClient.Do(forecastReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast data: %w", err)
	}
	defer forecastResp.Body.Close()

	if forecastResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(forecastResp.Body)
		return nil, fmt.Errorf("forecast API error: %s - %s", forecastResp.Status, string(body))
	}

	var forecast ForecastResponse
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}

	return &forecast, nil
}

// GetHourlyForecastByCoordinates gets hourly weather forecast for given latitude and longitude
// This can provide up to 156 hours (6.5 days) of hourly forecast data
func (w *WeatherClient) GetHourlyForecastByCoordinates(lat, lon float64) (*ForecastResponse, error) {
	// First, get the grid information for the coordinates
	pointsURL := fmt.Sprintf("%s/points/%.4f,%.4f", w.BaseURL, lat, lon)
	
	req, err := http.NewRequest("GET", pointsURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("User-Agent", "weather-app/1.0 (your-email@example.com)")
	
	resp, err := w.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get points data: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("NWS API error: %s - %s", resp.Status, string(body))
	}
	
	var pointsResp PointsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pointsResp); err != nil {
		return nil, fmt.Errorf("failed to decode points response: %w", err)
	}
	
	// Use the hourly forecast URL instead of the regular forecast URL
	forecastReq, err := http.NewRequest("GET", pointsResp.Properties.ForecastHourly, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create forecast request: %w", err)
	}
	
	forecastReq.Header.Set("User-Agent", "weather-app/1.0 (your-email@example.com)")
	
	forecastResp, err := w.HTTPClient.Do(forecastReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get forecast data: %w", err)
	}
	defer forecastResp.Body.Close()
	
	if forecastResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(forecastResp.Body)
		return nil, fmt.Errorf("forecast API error: %s - %s", forecastResp.Status, string(body))
	}
	
	var forecast ForecastResponse
	if err := json.NewDecoder(forecastResp.Body).Decode(&forecast); err != nil {
		return nil, fmt.Errorf("failed to decode forecast response: %w", err)
	}
	
	return &forecast, nil
}

// FormatForecast returns a formatted string representation of the forecast
func (f *ForecastResponse) FormatForecast(periods int) string {
	if periods <= 0 || periods > len(f.Properties.Periods) {
		periods = len(f.Properties.Periods)
	}

	result := "Weather Forecast:\n"
	result += "==================\n\n"

	for i := 0; i < periods && i < len(f.Properties.Periods); i++ {
		period := f.Properties.Periods[i]
		result += fmt.Sprintf("📅 %s\n", period.Name)
		result += fmt.Sprintf("🌡️  Temperature: %d°%s", period.Temperature, period.TemperatureUnit)
		if period.TemperatureTrend != "" {
			result += fmt.Sprintf(" (%s)", period.TemperatureTrend)
		}
		result += "\n"
		result += fmt.Sprintf("💨 Wind: %s %s\n", period.WindSpeed, period.WindDirection)
		result += fmt.Sprintf("☁️  Conditions: %s\n", period.ShortForecast)
		if period.DetailedForecast != "" {
			result += fmt.Sprintf("📝 Details: %s\n", period.DetailedForecast)
		}
		result += "\n"
	}

	return result
}
