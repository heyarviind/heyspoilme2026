package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"heyspoilme/internal/config"
	"heyspoilme/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: import-cities <path-to-worldcities.csv>")
	}

	csvPath := os.Args[1]

	// Load config and connect to database
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create cities table if it doesn't exist
	log.Println("Creating cities table if not exists...")
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS cities (
			id SERIAL PRIMARY KEY,
			city VARCHAR(200) NOT NULL,
			city_ascii VARCHAR(200) NOT NULL,
			state VARCHAR(200) NOT NULL,
			latitude DOUBLE PRECISION NOT NULL,
			longitude DOUBLE PRECISION NOT NULL,
			population BIGINT DEFAULT 0,
			country_code CHAR(2) NOT NULL DEFAULT 'IN'
		)
	`)
	if err != nil {
		log.Fatal("Failed to create cities table:", err)
	}

	// Create indexes
	log.Println("Creating indexes...")
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_cities_city_ascii ON cities(city_ascii)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_cities_country_code ON cities(country_code)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_cities_search ON cities(country_code, city_ascii)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS idx_cities_population ON cities(population DESC)`)

	// Open CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read header
	header, err := reader.Read()
	if err != nil {
		log.Fatal("Failed to read CSV header:", err)
	}

	// Map header to column indices
	colIdx := make(map[string]int)
	for i, col := range header {
		colIdx[col] = i
	}

	// Verify required columns exist
	required := []string{"city", "city_ascii", "lat", "lng", "admin_name", "iso2", "population"}
	for _, col := range required {
		if _, ok := colIdx[col]; !ok {
			log.Fatalf("Missing required column: %s", col)
		}
	}

	// Prepare insert statement
	stmt, err := db.Prepare(`
		INSERT INTO cities (city, city_ascii, state, latitude, longitude, population, country_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		log.Fatal("Failed to prepare statement:", err)
	}
	defer stmt.Close()

	// Clear existing cities
	log.Println("Clearing existing cities...")
	_, err = db.Exec("TRUNCATE TABLE cities RESTART IDENTITY")
	if err != nil {
		log.Fatal("Failed to truncate cities table:", err)
	}

	// Process rows
	count := 0
	errors := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading row: %v", err)
			errors++
			continue
		}

		city := record[colIdx["city"]]
		cityASCII := record[colIdx["city_ascii"]]
		state := record[colIdx["admin_name"]]
		latStr := record[colIdx["lat"]]
		lngStr := record[colIdx["lng"]]
		countryCode := record[colIdx["iso2"]]
		popStr := record[colIdx["population"]]

		// Parse latitude
		lat, err := strconv.ParseFloat(latStr, 64)
		if err != nil {
			log.Printf("Invalid latitude for %s: %v", city, err)
			errors++
			continue
		}

		// Parse longitude
		lng, err := strconv.ParseFloat(lngStr, 64)
		if err != nil {
			log.Printf("Invalid longitude for %s: %v", city, err)
			errors++
			continue
		}

		// Parse population (may be empty)
		var population int64 = 0
		if popStr != "" {
			population, _ = strconv.ParseInt(popStr, 10, 64)
		}

		// Normalize state name for India (remove diacritics common in source)
		if countryCode == "IN" {
			state = normalizeIndianState(state)
		}

		// Insert
		_, err = stmt.Exec(city, cityASCII, state, lat, lng, population, countryCode)
		if err != nil {
			log.Printf("Failed to insert %s: %v", city, err)
			errors++
			continue
		}

		count++
		if count%5000 == 0 {
			log.Printf("Imported %d cities...", count)
		}
	}

	log.Printf("Import complete! Inserted %d cities with %d errors", count, errors)
}

// normalizeIndianState converts state names with diacritics to standard names
func normalizeIndianState(state string) string {
	// Common mappings from CSV diacritical names to standard names
	mappings := map[string]string{
		"Mahārāshtra":        "Maharashtra",
		"Karnātaka":          "Karnataka",
		"Tamil Nādu":         "Tamil Nadu",
		"Gujarāt":            "Gujarat",
		"Rājasthān":          "Rajasthan",
		"Bihār":              "Bihar",
		"Andhra Pradesh":     "Andhra Pradesh",
		"Telangāna":          "Telangana",
		"Madhya Pradesh":     "Madhya Pradesh",
		"Uttar Pradesh":      "Uttar Pradesh",
		"West Bengal":        "West Bengal",
		"Kerala":             "Kerala",
		"Punjab":             "Punjab",
		"Haryāna":            "Haryana",
		"Jhārkhand":          "Jharkhand",
		"Chhattīsgarh":       "Chhattisgarh",
		"Assam":              "Assam",
		"Odisha":             "Odisha",
		"Uttarākhand":        "Uttarakhand",
		"Himāchal Pradesh":   "Himachal Pradesh",
		"Tripura":            "Tripura",
		"Meghālaya":          "Meghalaya",
		"Manipur":            "Manipur",
		"Nāgāland":           "Nagaland",
		"Goa":                "Goa",
		"Arunāchal Pradesh":  "Arunachal Pradesh",
		"Mizoram":            "Mizoram",
		"Sikkim":             "Sikkim",
		"Delhi":              "Delhi",
		"Jammu and Kashmīr":  "Jammu and Kashmir",
		"Jammu and Kashmir":  "Jammu and Kashmir",
		"Ladākh":             "Ladakh",
		"Puducherry":         "Puducherry",
		"Chandīgarh":         "Chandigarh",
		"Chandigarh":         "Chandigarh",
		"Andaman and Nicobar": "Andaman and Nicobar",
		"Dādra and Nagar Haveli and Damān and Diu": "Dadra and Nagar Haveli and Daman and Diu",
		"Lakshadweep":        "Lakshadweep",
	}

	// Try exact match first
	if normalized, ok := mappings[state]; ok {
		return normalized
	}

	// Try case-insensitive match
	stateLower := strings.ToLower(state)
	for k, v := range mappings {
		if strings.ToLower(k) == stateLower {
			return v
		}
	}

	// Return original if no mapping found
	return state
}

