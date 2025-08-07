# Weather CLI Tool

A command-line weather application that fetches forecast data from the National Weather Service API and stores it in a MySQL database for historical tracking.

## Features

- **Daily Forecasts**: Get traditional day/night period forecasts with detailed descriptions
- **Hourly Forecasts**: Get granular hour-by-hour weather data (up to 156 hours)
- **Database Storage**: Save forecast data to MySQL for historical tracking
- **Historical Data**: Retrieve previously saved forecast data
- **Duplicate Prevention**: Automatic deduplication of forecast records
- **Configuration Support**: Use config files or command-line flags
- **Coordinate Validation**: Input validation for latitude/longitude values

## Installation

1. Clone the repository
2. Install dependencies: `go mod tidy`
3. Build the application: `go build`
4. Configure your database connection in `.weather.yml`

## Usage

### Get Current Forecast
```bash
# Daily forecast
./weather forecast --lat 39.7391 --lon -104.9847

# Hourly forecast
./weather forecast --hourly --lat 39.7391 --lon -104.9847

# Save to database
./weather forecast --save --lat 39.7391 --lon -104.9847
```

### View Historical Data
```bash
# Daily historical forecasts
./weather history --lat 39.7391 --lon -104.9847

# Hourly historical forecasts
./weather history --hourly --lat 39.7391 --lon -104.9847

# Specify number of periods
./weather history --periods 10 --lat 39.7391 --lon -104.9847
```

## Automated Data Collection with Cron

For continuous weather data collection, you can set up cron jobs to automatically save forecast data at regular intervals. Below are recommended crontab entries for different use cases:

### Basic Automated Collection

```bash
# Edit your crontab
crontab -e

# Add these entries (adjust paths and coordinates as needed):

# Daily forecast collection - every 6 hours
0 */6 * * * /path/to/weather forecast --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-daily.log 2>&1

# Hourly forecast collection - every 2 hours
0 */2 * * * /path/to/weather forecast --hourly --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-hourly.log 2>&1
```

### Comprehensive Automated Collection

```bash
# Multiple location daily forecasts - every 4 hours
0 */4 * * * /path/to/weather forecast --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-denver.log 2>&1
0 */4 * * * /path/to/weather forecast --save --lat 40.7128 --lon -74.0060 >> /var/log/weather-nyc.log 2>&1
0 */4 * * * /path/to/weather forecast --save --lat 34.0522 --lon -118.2437 >> /var/log/weather-la.log 2>&1

# Hourly forecasts for primary location - every hour during active hours
0 6-22 * * * /path/to/weather forecast --hourly --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-hourly.log 2>&1

# Daily forecasts for primary location - twice daily
0 6,18 * * * /path/to/weather forecast --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-daily.log 2>&1
```

### High-Frequency Monitoring

```bash
# For weather-sensitive operations - collect every 30 minutes during business hours
*/30 6-18 * * 1-5 /path/to/weather forecast --hourly --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-business.log 2>&1

# Weekend data collection - every 2 hours
0 */2 * * 6,7 /path/to/weather forecast --hourly --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-weekend.log 2>&1
```

### Storm Season Enhanced Collection

```bash
# Enhanced collection during storm season (April-September)
# Hourly forecasts every 30 minutes
*/30 * * 4-9 * /path/to/weather forecast --hourly --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-storm-season.log 2>&1

# Daily forecasts every 2 hours during storm season
0 */2 * * 4-9 /path/to/weather forecast --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-storm-daily.log 2>&1

# Regular collection during off-season (October-March)
0 */4 * * 10-3 /path/to/weather forecast --save --lat 39.7391 --lon -104.9847 >> /var/log/weather-offseason.log 2>&1
```

### Crontab Schedule Explanation

| Schedule | Description |
|----------|-------------|
| `0 */6 * * *` | Every 6 hours at the top of the hour |
| `0 */2 * * *` | Every 2 hours at the top of the hour |
| `*/30 6-18 * * 1-5` | Every 30 minutes, 6 AM to 6 PM, Monday-Friday |
| `0 6,18 * * *` | Twice daily at 6 AM and 6 PM |
| `0 */2 * * 6,7` | Every 2 hours on weekends |
| `*/30 * * 4-9 *` | Every 30 minutes during April-September |

### Setup Instructions

1. **Update paths**: Replace `/path/to/weather` with the actual path to your weather binary
2. **Set coordinates**: Update latitude and longitude values for your location(s)
3. **Create log directories**: Ensure log directories exist and are writable
4. **Test commands**: Run the commands manually first to ensure they work
5. **Monitor logs**: Check log files regularly to ensure data collection is working

### Log Management

Consider adding log rotation to prevent log files from growing too large:

```bash
# Add to /etc/logrotate.d/weather
/var/log/weather*.log {
    weekly
    rotate 52
    compress
    delaycompress
    missingok
    notifempty
    copytruncate
}
```

### Database Maintenance

The application automatically handles duplicate prevention, but you may want to periodically clean old data:

```bash
# Monthly cleanup of data older than 1 year (add to crontab)
0 2 1 * * mysql -u username -p database_name -e "DELETE FROM weather_forecasts WHERE forecast_date < DATE_SUB(NOW(), INTERVAL 1 YEAR);" >> /var/log/weather-cleanup.log 2>&1
```

## Configuration

Create a `.weather.yml` file in the application directory:

```yaml
database:
  host: "localhost"
  port: 3306
  user: "weather_user"
  password: "your_password"
  dbname: "weather_db"
  
forecast:
  latitude: 39.7391
  longitude: -104.9847
  periods: 7
  save: false
  hourly: false
```

## Database Schema

The application uses a `weather_forecasts` table with the following structure:
- Location coordinates (latitude, longitude)
- Forecast metadata (date, period number, forecast type)
- Weather data (temperature, wind, conditions, etc.)
- Temporal data (start/end times)
- Forecast type indicator (daily vs hourly)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
