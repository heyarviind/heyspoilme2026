package models

type City struct {
	ID          int     `json:"id" db:"id"`
	City        string  `json:"city" db:"city"`
	CityASCII   string  `json:"city_ascii" db:"city_ascii"`
	State       string  `json:"state" db:"state"`
	Latitude    float64 `json:"latitude" db:"latitude"`
	Longitude   float64 `json:"longitude" db:"longitude"`
	Population  int64   `json:"population" db:"population"`
	CountryCode string  `json:"country_code" db:"country_code"`
}

type CitySearchResult struct {
	City      string  `json:"city"`
	State     string  `json:"state"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

