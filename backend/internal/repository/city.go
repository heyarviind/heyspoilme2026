package repository

import (
	"database/sql"

	"heyspoilme/internal/models"
)

type CityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) *CityRepository {
	return &CityRepository{db: db}
}

// SearchCities searches for cities by name prefix (case-insensitive)
// Returns up to 10 results, ordered by population (most populated first)
func (r *CityRepository) SearchCities(query string, countryCode string) ([]models.CitySearchResult, error) {
	if len(query) < 2 {
		return []models.CitySearchResult{}, nil
	}

	rows, err := r.db.Query(`
		SELECT city, state, latitude, longitude
		FROM cities
		WHERE country_code = $1 
		  AND city_ascii ILIKE $2
		ORDER BY population DESC NULLS LAST
		LIMIT 10
	`, countryCode, query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CitySearchResult
	for rows.Next() {
		var c models.CitySearchResult
		if err := rows.Scan(&c.City, &c.State, &c.Latitude, &c.Longitude); err != nil {
			return nil, err
		}
		results = append(results, c)
	}

	return results, nil
}

// GetCityByName gets a city by exact name match (case-insensitive)
func (r *CityRepository) GetCityByName(cityName string, countryCode string) (*models.CitySearchResult, error) {
	var c models.CitySearchResult
	err := r.db.QueryRow(`
		SELECT city, state, latitude, longitude
		FROM cities
		WHERE country_code = $1 
		  AND (city_ascii ILIKE $2 OR city ILIKE $2)
		ORDER BY population DESC NULLS LAST
		LIMIT 1
	`, countryCode, cityName).Scan(&c.City, &c.State, &c.Latitude, &c.Longitude)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &c, nil
}

