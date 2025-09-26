# hourstable

[![Build Status](https://travis-ci.org/geniusrabbit/hourstable.svg?branch=master)](https://travis-ci.org/geniusrabbit/hourstable)
[![Go Report Card](https://goreportcard.com/badge/github.com/geniusrabbit/hourstable)](https://goreportcard.com/report/github.com/geniusrabbit/hourstable)
[![GoDoc](https://godoc.org/github.com/geniusrabbit/hourstable?status.svg)](https://godoc.org/github.com/geniusrabbit/hourstable)
[![Coverage Status](https://coveralls.io/repos/github/geniusrabbit/hourstable/badge.svg)](https://coveralls.io/github/geniusrabbit/hourstable)

A Go package that provides efficient data types for representing weekly hour schedules. Perfect for managing business hours, availability windows, scheduling systems, and time-based configurations.

## Features

- 🕒 **Weekly Hour Management**: Efficiently represent active/inactive hours for each day of the week
- 💾 **Database Integration**: Built-in SQL database support with `driver.Valuer` and `sql.Scanner` interfaces
- 🔄 **Multiple Serialization Formats**: Support for string, JSON, YAML, and structured formats
- ⚡ **Memory Efficient**: Bit-packed internal representation for optimal performance
- 🧪 **Well Tested**: Comprehensive test suite with benchmarks
- 📦 **Zero Dependencies**: Pure Go implementation with only standard library dependencies

## Installation

```bash
go get github.com/geniusrabbit/hourstable
```

**For YAML support:**

```bash  
go get gopkg.in/yaml.v3
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    
    "github.com/geniusrabbit/hourstable"
)

func main() {
    // Create business hours: Mon-Fri 9AM-5PM
    // String format: 7 days × 24 hours = 168 characters
    // Day order: Sunday, Monday, Tuesday, Wednesday, Thursday, Friday, Saturday
    // Hour format: 0-23 (24-hour format)
    businessHours := hourstable.MustHoursByString(
        "000000000111111111000000" + // Sunday: hours 0-8 off, 9-17 on, 18-23 off
        "000000000111111111000000" + // Monday: 9AM-5PM (hours 9-17)
        "000000000111111111000000" + // Tuesday: 9AM-5PM
        "000000000111111111000000" + // Wednesday: 9AM-5PM
        "000000000111111111000000" + // Thursday: 9AM-5PM
        "000000000111111111000000" + // Friday: 9AM-5PM
        "000000000000000000000000",  // Saturday: closed all day
    )
    
    // Check if current time is within business hours
    fmt.Printf("Open now? %v\n", businessHours.TestTime(time.Now()))
    
    // Check specific day and hour (using 24-hour format: 0-23)
    fmt.Printf("Open Monday 10AM? %v\n", businessHours.TestHour(time.Monday, 10))
}
```

## Hours Table Example

| Day of week / Hour | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 | 15 | 16 | 17 | 18 | 19 | 20 | 21 | 22 | 23 |
|:-------------------|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|:--:|
| Sunday             | X | X | X |   |   |   |   |   |   |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Monday             | X | X | X | X | X | X | X | X | X |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Tuesday            | X | X | X | X | X | X | X | X | X |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Wednesday          | X | X | X | X | X | X | X | X | X |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Thursday           | X | X | X | X | X | X | X | X | X |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Friday             | X | X | X | X | X | X | X | X | X |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |
| Saturday           | X | X | X |   |   |   |   |   |   |   |    |    |    |    |    |    |    |    |    |    |    |    | X  | X  |

In a database, it could be stored as TEXT of '1' and '0' symbols

```plain
111111111111111111111111111111111111111111111111000000000000000000000000000000000000000000000000000000000000000000000000111111111111111111111111111111111111111111111111
```

or as the JSON object

```json
{
  "mon": "111111111111111111111111",
  "tue": "111111111111111111111111",
  "wed": "000000000000000000000000",
  "thu": "000000000000000000000000",
  "fri": "000000000000000000000000",
  "sat": "111111111111111111111111",
  "sun": "000000001111111111111111"
}
```

or short form of JSON with wildcards:

```json
{
  "mon": "*",
  "tue": "*", 
  "sat": "*",
  "sun": "000000001111111111111111"
}
```

**Format Explanations:**

- `"*"` - All 24 hours are active (shorthand)
- `""` - All 24 hours are inactive (empty string)
- `"1"` - Hour is active
- `"0"` - Hour is inactive

**Important Notes:**

- **Day Order**: String format follows `time.Weekday` order: Sunday(0), Monday(1), Tuesday(2), Wednesday(3), Thursday(4), Friday(5), Saturday(6)
- **Hour Format**: Uses 24-hour format (0-23), where hour 0 = midnight, hour 12 = noon
- **String Length**: Full format requires exactly 168 characters (7 days × 24 hours)
- **Timezone**: The package doesn't handle timezones - all times are treated as local time
- **Internal Storage**: Uses bit-packed representation for memory efficiency (24 bytes per schedule)

## API Reference

### Core Types

#### Hours

```go
type Hours []byte
```

Main type for representing weekly hour schedules.

#### HoursObject

```go
type HoursObject Hours
```

Alternative type with JSON-optimized serialization.

### Creation Functions

```go
// Create from string format
func HoursByString(s string) (Hours, error)
func MustHoursByString(s string) Hours  // Panics on error

// Create from JSON format  
func HoursByJSON(data []byte) (Hours, error)
```

### Query Methods

```go
// Test specific hour
func (h Hours) TestHour(weekDay time.Weekday, hour byte) bool
func (h Hours) TestTime(t time.Time) bool

// Check schedule state
func (h Hours) IsAllActive() bool
func (h Hours) IsNoActive() bool
func (h Hours) Equal(h2 Hours) bool
```

### Modification Methods

```go
// Set specific hour active/inactive
func (h *Hours) SetHour(weekDay time.Weekday, hour byte, active bool)

// Merge schedules
func (h Hours) Merge(h2 Hours)

// Create copy
func (h Hours) Clone() Hours
```

### Serialization

```go
// String representation
func (h Hours) String() string

// JSON serialization
func (h Hours) MarshalJSON() ([]byte, error)
func (h *Hours) UnmarshalJSON(data []byte) error

// YAML serialization
func (h Hours) MarshalYAML() (any, error)
func (h *Hours) UnmarshalYAML(node *yaml.Node) error

// Database integration
func (h Hours) Value() (driver.Value, error)  // driver.Valuer
func (h *Hours) Scan(value any) error         // sql.Scanner
```

## Usage Examples

### Business Hours Management

```go
// Restaurant hours: Open 11AM-10PM daily, closed Mondays
restaurant := make(hourstable.Hours, 24) // Initialize with 24 hours
for day := time.Tuesday; day <= time.Sunday; day++ {
    for hour := 11; hour <= 21; hour++ { // 11AM-9PM (hour 21 = 9PM)
        restaurant.SetHour(day, byte(hour), true)
    }
}

// Check if open at specific time
isOpen := restaurant.TestTime(time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC))
fmt.Printf("Restaurant open at 2:30 PM on Jan 15? %v\n", isOpen)
```

### Database Storage

```go
import "database/sql"

// Store in database
_, err := db.Exec("INSERT INTO venues (name, hours) VALUES (?, ?)", 
                  "Coffee Shop", businessHours)

// Read from database  
var hours hourstable.Hours
err = db.QueryRow("SELECT hours FROM venues WHERE id = ?", 1).Scan(&hours)
```

### JSON API Integration

```go
// HTTP handler example
func updateSchedule(w http.ResponseWriter, r *http.Request) {
    var schedule hourstable.Hours
    if err := json.NewDecoder(r.Body).Decode(&schedule); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Use schedule...
    json.NewEncoder(w).Encode(schedule)
}
```

### YAML Configuration Files

```go
import (
    "os"
    "gopkg.in/yaml.v3"
    "github.com/geniusrabbit/hourstable"
)

// Configuration file example
type Config struct {
    BusinessHours hourstable.HoursObject `yaml:"business_hours"`
    MaintenanceWindow hourstable.Hours `yaml:"maintenance_window"`
}

// Load from YAML file
func loadConfig(filename string) (*Config, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var config Config
    err = yaml.Unmarshal(data, &config)
    return &config, err
}

// Save to YAML file  
func saveConfig(config *Config, filename string) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }
    return os.WriteFile(filename, data, 0644)
}
```

**Example YAML configuration:**

```yaml
# Simple string format
business_hours: "000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000111111110000000000000000000000000000000"

# Structured format (using HoursObject)
maintenance_window:
  mon: "000000001100000000000000"  # 1-2 AM
  tue: "000000001100000000000000"
  wed: "000000001100000000000000" 
  thu: "000000001100000000000000"
  fri: "000000001100000000000000"
  sat: ""  # No maintenance
  sun: ""  # No maintenance
```

### Helper Functions

```go
// Create schedule for specific hour range
func createDailySchedule(startHour, endHour int) hourstable.Hours {
    hours := make(hourstable.Hours, 24)
    for day := time.Sunday; day <= time.Saturday; day++ {
        for hour := startHour; hour < endHour; hour++ {
            hours.SetHour(day, byte(hour), true)
        }
    }
    return hours
}
```

## Use Cases

- **Business Hours**: Store and validate operating hours for businesses
- **Scheduling Systems**: Define availability windows for resources or staff  
- **Content Management**: Control when content should be active/visible
- **Rate Limiting**: Define time-based access patterns
- **Automation**: Schedule when automated processes should run
- **Booking Systems**: Manage availability for appointments or reservations

## Performance

The package is optimized for performance with:

- Bit-packed storage using only 24 bytes per schedule
- O(1) hour lookup operations
- Efficient serialization/deserialization
- Special optimizations for common patterns (all active, none active)

## Testing & Benchmarks

```bash
go test -timeout 30s github.com/geniusrabbit/hourstable -v -race
```

### Benchmarks

```bash
go test -benchmem -run=^$ github.com/geniusrabbit/hourstable -bench . -v

# Example output:
# goos: darwin
# goarch: arm64
# pkg: github.com/geniusrabbit/hourstable
# cpu: Apple M2 Ultra
# Benchmark_Hours-24    651210    1766 ns/op    0 B/op    0 allocs/op
# PASS
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

## Author

- **Dmitry Ponomarev** - [@demdxx](https://github.com/demdxx)

## Related Projects

- Looking for more scheduling utilities? Check out other [GeniusRabbit](https://github.com/geniusrabbit) projects!
