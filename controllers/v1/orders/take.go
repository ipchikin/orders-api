package orders

import (
	"errors"
	"net/http"
	"orders-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TakeRequest
type TakeRequest struct {
	Status string `binding:"required"`
}

// TakeOrder controller
func TakeOrder(c *gin.Context) {
	var takeRequest TakeRequest
	err := c.BindJSON(&takeRequest)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid request body"))
	}

	// Check if TakeRequest status is valid
	if takeRequest.Status != "TAKEN" {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid request status value"))
	}

	// Check if id is a valid uuid
	id := c.Param("id")
	_, err = uuid.Parse(id)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Invalid request id value"))
	}

	// Take order
	om := models.OrdersModel{BaseModel: c.MustGet("db").(models.BaseModel)}
	err = om.Take(id, takeRequest.Status)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("Fail to take the order"))
	}

	c.JSON(http.StatusOK, &gin.H{"status": "SUCCESS"})
}
