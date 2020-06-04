package orders

import (
	"encoding/json"
	"errors"
	"net/http"
	"orders-api/configs"
	"orders-api/models"
	"regexp"

	"github.com/gin-gonic/gin"
)

// PlaceData
type PlaceData struct {
	Origin      [2]string `json:"origin" binding:"required"`
	Destination [2]string `json:"destination" binding:"required"`
}

// DistanceMatrixAPIResponse
type DistanceMatrixAPIResponse struct {
	Rows   []Row
	Status string
}

type Row struct {
	Elements []Element
}

type Element struct {
	Distance Distance
	Status   string
}

type Distance struct {
	Value uint32
}

// PlaceOrder
func PlaceOrder(c *gin.Context) {
	var placeData PlaceData
	err := c.BindJSON(&placeData)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	// Validate origin coordinates
	ok := validateLatitude(placeData.Origin[0])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid origin latitude"))
	}
	ok = validateLongitude(placeData.Origin[1])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid origin longitude"))
	}

	// Validate destination coordinates
	ok = validateLatitude(placeData.Destination[0])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid destination latitude"))
	}
	ok = validateLongitude(placeData.Destination[1])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid destination longitude"))
	}

	// Get config from gin context
	cfg := c.MustGet("config").(configs.Config)

	// Call Google Distance Matrix API
	req, err := http.NewRequest(cfg.DistanceMatrixAPI.Method, cfg.DistanceMatrixAPI.URL, nil)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	q := req.URL.Query()
	q.Add("origins", placeData.Origin[0]+","+placeData.Origin[1])
	q.Add("destinations", placeData.Destination[0]+","+placeData.Destination[1])
	q.Add("key", cfg.DistanceMatrixAPI.Key)
	req.URL.RawQuery = q.Encode()

	// Get http client from gin context
	client := c.MustGet("client").(*http.Client)
	resp, err := client.Do(req)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Fail to get distance details"))
	}

	// Parse response body to struct
	apiResp := DistanceMatrixAPIResponse{}
	json.NewDecoder(resp.Body).Decode(&apiResp)

	// Abort if status not ok
	if apiResp.Status != "OK" {
		c.AbortWithError(http.StatusInternalServerError, errors.New("Fail to get distance details"))
	}

	if apiResp.Rows[0].Elements[0].Status != "OK" {
		c.AbortWithError(http.StatusBadRequest, errors.New("No route could be found between the origin and destination"))
	}

	distance := apiResp.Rows[0].Elements[0].Distance.Value
	// Check if distance is valid
	if distance == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New("Distance between the origin and destination is too small"))
	}

	status := "UNASSIGNED"

	// Use orders model
	om := models.OrdersModel{BaseModel: c.MustGet("db").(models.BaseModel)}
	id, err := om.Place(placeData.Origin, placeData.Destination, distance, status)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

	c.JSON(http.StatusOK, &gin.H{
		"id":       id,
		"distance": distance,
		"status":   status,
	})
}

// validateLatitude
func validateLatitude(lat string) bool {
	matched, err := regexp.MatchString(`^(\+|-)?(?:90(?:(?:\.0+)?)|(?:[0-9]|[1-8][0-9])(?:(?:\.[0-9]+)?))$`, lat)
	if err != nil {
		return false
	}

	return matched
}

// validateLongitude
func validateLongitude(long string) bool {
	matched, err := regexp.MatchString(`^(\+|-)?(?:180(?:(?:\.0+)?)|(?:[0-9]|[1-9][0-9]|1[0-7][0-9])(?:(?:\.[0-9]+)?))$`, long)
	if err != nil {
		return false
	}

	return matched
}
