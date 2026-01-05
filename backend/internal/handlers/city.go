package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"heyspoilme/internal/repository"
)

type CityHandler struct {
	cityRepo *repository.CityRepository
}

func NewCityHandler(cityRepo *repository.CityRepository) *CityHandler {
	return &CityHandler{cityRepo: cityRepo}
}

// SearchCities handles GET /api/cities/search?q=<query>&country=<country_code>
func (h *CityHandler) SearchCities(c *gin.Context) {
	query := c.Query("q")
	if len(query) < 2 {
		c.JSON(http.StatusOK, gin.H{"cities": []any{}})
		return
	}

	// Default to India, but allow other countries
	countryCode := c.DefaultQuery("country", "IN")

	cities, err := h.cityRepo.SearchCities(query, countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search cities"})
		return
	}

	if cities == nil {
		c.JSON(http.StatusOK, gin.H{"cities": []any{}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"cities": cities})
}

