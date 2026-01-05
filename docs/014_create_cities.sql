-- Create cities table for autocomplete
-- To apply: Run this in PocketBase admin SQL editor or psql

CREATE TABLE IF NOT EXISTS cities (
    id SERIAL PRIMARY KEY,
    city VARCHAR(200) NOT NULL,
    city_ascii VARCHAR(200) NOT NULL,
    state VARCHAR(200) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    population BIGINT DEFAULT 0,
    country_code CHAR(2) NOT NULL DEFAULT 'IN'
);

-- Create indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_cities_city_ascii ON cities(city_ascii);
CREATE INDEX IF NOT EXISTS idx_cities_country_code ON cities(country_code);
CREATE INDEX IF NOT EXISTS idx_cities_search ON cities(country_code, city_ascii);
CREATE INDEX IF NOT EXISTS idx_cities_population ON cities(population DESC);

-- For text search
CREATE INDEX IF NOT EXISTS idx_cities_city_trgm ON cities USING gin (city_ascii gin_trgm_ops);

-- Enable pg_trgm extension if not already enabled
CREATE EXTENSION IF NOT EXISTS pg_trgm;

