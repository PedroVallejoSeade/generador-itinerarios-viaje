// Package city loads the bundled world-cities dataset and provides city search.
package city

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

//go:embed cities.csv
var citiesCSV []byte

// City is a populated place a user can search for and disambiguate.
type City struct {
	Name       string
	Country    string
	Region     string
	Population int64
}

// Load parses the bundled dataset embedded at build time and returns the cities.
func Load() ([]City, error) {
	return parse(bytes.NewReader(citiesCSV))
}

// parse reads CSV rows (city,country,admin_name,population) from r. It skips the
// header row and any row with an empty city name, and defaults a blank or
// non-numeric population to 0. A malformed CSV stream returns an error (FR-009).
func parse(r io.Reader) ([]City, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	var cities []City
	first := true
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parse city data: %w", err)
		}

		// Skip the header row.
		if first {
			first = false
			if len(record) > 0 && strings.EqualFold(strings.TrimSpace(record[0]), "city") {
				continue
			}
		}

		if len(record) < 4 {
			continue
		}

		name := strings.TrimSpace(record[0])
		if name == "" {
			continue
		}

		population, _ := strconv.ParseInt(strings.TrimSpace(record[3]), 10, 64)

		cities = append(cities, City{
			Name:       name,
			Country:    strings.TrimSpace(record[1]),
			Region:     strings.TrimSpace(record[2]),
			Population: population,
		})
	}

	return cities, nil
}
