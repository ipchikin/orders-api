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

// PlaceRequest
type PlaceRequest struct {
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
	var placeRequest PlaceRequest
	err := c.BindJSON(&placeRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	// Validate origin coordinates
	ok := validateLatitude(placeRequest.Origin[0])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid origin latitude"))
	}
	ok = validateLongitude(placeRequest.Origin[1])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid origin longitude"))
	}

	// Validate destination coordinates
	ok = validateLatitude(placeRequest.Destination[0])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid destination latitude"))
	}
	ok = validateLongitude(placeRequest.Destination[1])
	if !ok {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid destination longitude"))
	}

	// Get the distance between origin and destination
	distance, err := callDistanceMatrixAPI(c, placeRequest.Origin, placeRequest.Destination)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Fail to get the distance between origin and destination"))
	}

	status := "UNASSIGNED"

	// Place order
	om := models.OrdersModel{BaseModel: c.MustGet("db").(models.BaseModel)}
	id, err := om.Place(placeRequest.Origin, placeRequest.Destination, distance, status)
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

// callDistanceMatrixAPI
func callDistanceMatrixAPI(c *gin.Context, origin, destination [2]string) (distance uint32, err error) {
	// Get config from gin context
	cfg := c.MustGet("config").(configs.Config)

	// Call Google Distance Matrix API
	req, err := http.NewRequest(cfg.DistanceMatrixAPI.Method, cfg.DistanceMatrixAPI.URL, nil)
	if err != nil {
		return
	}

	q := req.URL.Query()
	q.Add("origins", origin[0]+","+origin[1])
	q.Add("destinations", destination[0]+","+destination[1])
	q.Add("key", cfg.DistanceMatrixAPI.Key)
	req.URL.RawQuery = q.Encode()

	// Get http client from gin context
	client := c.MustGet("client").(*http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.New("Fail to get the distance between origin and destination")
		return
	}

	// Parse response body to struct
	apiResp := DistanceMatrixAPIResponse{}
	json.NewDecoder(resp.Body).Decode(&apiResp)

	// Abort if status not ok
	if apiResp.Status != "OK" {
		err = errors.New("Fail to get the distance between origin and destination")
		return
	}

	if apiResp.Rows[0].Elements[0].Status != "OK" {
		err = errors.New("Fail to get the distance between origin and destination")
		return
	}

	distance = apiResp.Rows[0].Elements[0].Distance.Value
	// Check if distance is valid
	if distance == 0 {
		err = errors.New("Fail to get the distance between origin and destination")
	}

	return
}
