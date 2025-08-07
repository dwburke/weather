package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dwburke/weather/types"
)

var (
	forecastPeriods int
	latitude        float64
	longitude       float64
	saveToDb        bool
)

func init() {
	rootCmd.AddCommand(forecast)

	// Add flags for coordinates
	forecast.Flags().Float64VarP(&latitude, "lat", "a", 0.0, "Latitude for weather forecast")
	forecast.Flags().Float64VarP(&longitude, "lon", "o", 0.0, "Longitude for weather forecast")
	forecast.Flags().IntVarP(&forecastPeriods, "periods", "p", 7, "Number of forecast periods to show (each day has day/night periods)")
	forecast.Flags().BoolVarP(&saveToDb, "save", "s", false, "Save forecast data to database")

	// Bind flags to viper for configuration file support
	viper.BindPFlag("forecast.latitude", forecast.Flags().Lookup("lat"))
	viper.BindPFlag("forecast.longitude", forecast.Flags().Lookup("lon"))
	viper.BindPFlag("forecast.periods", forecast.Flags().Lookup("periods"))
	viper.BindPFlag("forecast.save", forecast.Flags().Lookup("save"))
}

var forecast = &cobra.Command{
	Use:   "forecast",
	Short: "Get weather forecast for a location",
	Long: `Get weather forecast from the National Weather Service (api.weather.gov) for specified coordinates.
	
Note: The NWS API returns forecast "periods" rather than full days. 
Each day typically has 2 periods: daytime and nighttime.
So requesting 6 periods gives you approximately 3 full days of forecast.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get coordinates from flags or config
		lat := viper.GetFloat64("forecast.latitude")
		lon := viper.GetFloat64("forecast.longitude")
		periods := viper.GetInt("forecast.periods")
		save := viper.GetBool("forecast.save")

		// Fallback to old config key if new one doesn't exist
		if periods == 0 {
			periods = viper.GetInt("forecast.days")
		}
		if periods == 0 {
			periods = 7 // default
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

		fmt.Printf("Getting weather forecast for coordinates: %.4f, %.4f\n", lat, lon)
		fmt.Printf("Showing %d forecast periods (each day typically has day/night periods)\n\n", periods)

		// Create weather client and get forecast
		client := types.NewWeatherClient()
		forecast, err := client.GetForecastByCoordinates(lat, lon)
		if err != nil {
			return fmt.Errorf("failed to get weather forecast: %w", err)
		}

		// Save to database if requested
		if save {
			fmt.Printf("Saving forecast data to database...\n")
			if err := types.SaveForecastToDB(forecast, lat, lon); err != nil {
				return fmt.Errorf("failed to save forecast to database: %w", err)
			}
			fmt.Printf("✅ Forecast data saved successfully!\n\n")
		}

		// Display the forecast
		fmt.Print(forecast.FormatForecast(periods))

		return nil
	},
}
