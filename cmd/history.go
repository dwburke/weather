package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	
	"github.com/dwburke/weather/types"
)

var (
	historyPeriods int
	historyLat     float64
	historyLon     float64
	historyHourly  bool
)

func init() {
	rootCmd.AddCommand(history)
	
	// Add flags for coordinates
	history.Flags().Float64VarP(&historyLat, "lat", "a", 0.0, "Latitude for weather history")
	history.Flags().Float64VarP(&historyLon, "lon", "o", 0.0, "Longitude for weather history")
	history.Flags().IntVarP(&historyPeriods, "periods", "p", 7, "Number of historical forecast periods to show")
	history.Flags().BoolVarP(&historyHourly, "hourly", "H", false, "Get hourly historical forecast instead of daily periods")
	
	// Bind flags to viper for configuration file support
	viper.BindPFlag("history.latitude", history.Flags().Lookup("lat"))
	viper.BindPFlag("history.longitude", history.Flags().Lookup("lon"))
	viper.BindPFlag("history.periods", history.Flags().Lookup("periods"))
	viper.BindPFlag("history.hourly", history.Flags().Lookup("hourly"))
}

var history = &cobra.Command{
	Use:   "history",
	Short: "Get historical weather forecast data from database",
	Long:  `Retrieve previously saved weather forecast data from the database for specified coordinates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get coordinates from flags or config, fallback to forecast config
		lat := viper.GetFloat64("history.latitude")
		lon := viper.GetFloat64("history.longitude")
		periods := viper.GetInt("history.periods")
		hourly := viper.GetBool("history.hourly")
		
		// Fallback to forecast coordinates if history coordinates not set
		if lat == 0.0 && lon == 0.0 {
			lat = viper.GetFloat64("forecast.latitude")
			lon = viper.GetFloat64("forecast.longitude")
		}
		
		if periods == 0 {
			if hourly {
				periods = 24 // default to 24 hours for hourly history
			} else {
				periods = 7 // default to 7 periods for daily history
			}
		}
		
		// Check if coordinates are provided
		if lat == 0.0 && lon == 0.0 {
			return fmt.Errorf("latitude and longitude must be provided. Use --lat and --lon flags or set them in config file")
		}
		
		if lat < -90 || lat > 90 {
			return fmt.Errorf("latitude must be between -90 and 90 degrees")
		}
		
		if lon < -180 || lon > 180 {
			return fmt.Errorf("longitude must be between -180 and 180 degrees")
		}
		
		forecastType := "daily"
		if hourly {
			forecastType = "hourly"
		}
		
		fmt.Printf("Getting historical weather forecast for coordinates: %.4f, %.4f\n", lat, lon)
		fmt.Printf("Showing %d historical %s forecast periods\n\n", periods, forecastType)
		
		// Get historical forecast data from database
		forecasts, err := types.GetLatestForecast(lat, lon, periods, hourly)
		if err != nil {
			return fmt.Errorf("failed to get historical forecast: %w", err)
		}
		
		if len(forecasts) == 0 {
			fmt.Printf("No historical %s forecast data found for coordinates %.4f, %.4f\n", forecastType, lat, lon)
			fmt.Printf("Use 'weather forecast --save' to save forecast data to the database first.\n")
			return nil
		}
		
		// Display the historical forecast
		fmt.Printf("Historical Weather Forecast (%s, saved: %s):\n", forecastType, forecasts[0].ForecastDate.Format("2006-01-02 15:04:05"))
		fmt.Printf("=========================================================\n\n")
		
		for _, forecast := range forecasts {
			fmt.Printf("📅 %s\n", forecast.Name)
			fmt.Printf("🌡️  Temperature: %d°%s", forecast.Temperature, forecast.TemperatureUnit)
			if forecast.TemperatureTrend != "" {
				fmt.Printf(" (%s)", forecast.TemperatureTrend)
			}
			fmt.Printf("\n")
			fmt.Printf("💨 Wind: %s %s\n", forecast.WindSpeed, forecast.WindDirection)
			fmt.Printf("☁️  Conditions: %s\n", forecast.ShortForecast)
			if forecast.DetailedForecast != "" {
				fmt.Printf("📝 Details: %s\n", forecast.DetailedForecast)
			}
			fmt.Printf("⏰ Period: %s to %s\n", 
				forecast.StartTime.Format("Jan 2 3:04 PM"), 
				forecast.EndTime.Format("Jan 2 3:04 PM"))
			fmt.Printf("\n")
		}
		
		return nil
	},
}
