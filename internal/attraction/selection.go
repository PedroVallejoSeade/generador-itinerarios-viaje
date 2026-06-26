package attraction

import "github.com/PedroVallejoSeade/generador-itinerarios-viaje/internal/city"

// CitySelection is the city the user chose from the search results, carrying
// enough identity to resolve the correct city for the attractions lookup
// (FR-001, FR-002a). It is transient — one selection→attractions exchange.
type CitySelection struct {
	Index   int    // 1-based position the user typed against the displayed list.
	Name    string // May contain accents/non-ASCII (e.g. "São Paulo").
	Country string // Disambiguates same-named cities.
	Region  string // Further disambiguation when available (may be empty).
}

// NewCitySelection maps a chosen city.City and its 1-based display index to a
// CitySelection used to resolve attractions (data-model.md).
func NewCitySelection(index int, c city.City) CitySelection {
	return CitySelection{
		Index:   index,
		Name:    c.Name,
		Country: c.Country,
		Region:  c.Region,
	}
}
