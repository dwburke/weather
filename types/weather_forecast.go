package types

import (
	"time"
	
	"github.com/dwburke/weather/db"
)

// WeatherForecast represents a weather forecast record in the database
type WeatherForecast struct {
	ID               uint      `json:"id" gorm:"primarykey"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	
	// Location information
	Latitude         float64   `json:"latitude" gorm:"column:latitude;not null"`
	Longitude        float64   `json:"longitude" gorm:"column:longitude;not null"`
	
	// Forecast period information
	PeriodNumber     int       `json:"period_number" gorm:"column:period_number;not null"`
	Name             string    `json:"name" gorm:"column:name;not null"`
	StartTime        time.Time `json:"start_time" gorm:"column:start_time;not null"`
	EndTime          time.Time `json:"end_time" gorm:"column:end_time;not null"`
	IsDaytime        bool      `json:"is_daytime" gorm:"column:is_daytime"`
	
	// Weather data
	Temperature      int       `json:"temperature" gorm:"column:temperature"`
	TemperatureUnit  string    `json:"temperature_unit" gorm:"column:temperature_unit"`
	TemperatureTrend string    `json:"temperature_trend" gorm:"column:temperature_trend"`
	WindSpeed        string    `json:"wind_speed" gorm:"column:wind_speed"`
	WindDirection    string    `json:"wind_direction" gorm:"column:wind_direction"`
	Icon             string    `json:"icon" gorm:"column:icon"`
	ShortForecast    string    `json:"short_forecast" gorm:"column:short_forecast"`
	DetailedForecast string    `json:"detailed_forecast" gorm:"column:detailed_forecast;type:text"`
	
	// Metadata
	ForecastDate     time.Time `json:"forecast_date" gorm:"column:forecast_date;index"` // When this forecast was retrieved
}

func (WeatherForecast) TableName() string {
	return "weather_forecasts"
}

// Create saves a new weather forecast record to the database
func (w *WeatherForecast) Create() error {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return err
	}

	if err := gdbh.Create(&w).Error; err != nil {
		return err
	}

	return nil
}

// Save updates an existing weather forecast record
func (w *WeatherForecast) Save() error {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return err
	}

	if err := gdbh.Save(&w).Error; err != nil {
		return err
	}

	return nil
}

// SaveForecastToDB saves a complete forecast response to the database
func SaveForecastToDB(forecast *ForecastResponse, lat, lon float64) error {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return err
	}

	// Auto-migrate the table if it doesn't exist
	gdbh.AutoMigrate(&WeatherForecast{})

	forecastDate := time.Now()
	
	// Save each forecast period
	for _, period := range forecast.Properties.Periods {
		startTime, err := time.Parse(time.RFC3339, period.StartTime)
		if err != nil {
			return err
		}
		
		endTime, err := time.Parse(time.RFC3339, period.EndTime)
		if err != nil {
			return err
		}
		
		weatherForecast := WeatherForecast{
			Latitude:         lat,
			Longitude:        lon,
			PeriodNumber:     period.Number,
			Name:             period.Name,
			StartTime:        startTime,
			EndTime:          endTime,
			IsDaytime:        period.IsDaytime,
			Temperature:      period.Temperature,
			TemperatureUnit:  period.TemperatureUnit,
			TemperatureTrend: period.TemperatureTrend,
			WindSpeed:        period.WindSpeed,
			WindDirection:    period.WindDirection,
			Icon:             period.Icon,
			ShortForecast:    period.ShortForecast,
			DetailedForecast: period.DetailedForecast,
			ForecastDate:     forecastDate,
		}
		
		if err := weatherForecast.Create(); err != nil {
			return err
		}
	}
	
	return nil
}

// GetLatestForecast retrieves the most recent forecast for given coordinates
func GetLatestForecast(lat, lon float64, limit int) ([]WeatherForecast, error) {
	gdbh, err := db.GetDB().DB()
	if err != nil {
		return nil, err
	}

	var forecasts []WeatherForecast
	
	// Get the most recent forecast date for these coordinates
	var latestDate time.Time
	if err := gdbh.Model(&WeatherForecast{}).
		Where("latitude = ? AND longitude = ?", lat, lon).
		Select("MAX(forecast_date)").
		Row().Scan(&latestDate); err != nil {
		return nil, err
	}
	
	// Get all periods from that forecast date
	query := gdbh.Where("latitude = ? AND longitude = ? AND forecast_date = ?", lat, lon, latestDate).
		Order("period_number ASC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	if err := query.Find(&forecasts).Error; err != nil {
		return nil, err
	}
	
	return forecasts, nil
}
