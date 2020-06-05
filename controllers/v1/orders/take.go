package orders

import (
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
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Check if TakeRequest status is valid
	if takeRequest.Status != "TAKEN" {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid request status value")
		return
	}

	// Check if id is a valid uuid
	id := c.Param("id")
	_, err = uuid.Parse(id)
	if err != nil {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid request id value")
		return
	}

	// Take order
	om := models.OrdersModel{BaseModel: c.MustGet("db").(models.BaseModel)}
	err = om.Take(id, takeRequest.Status)
	if err != nil {
		abortWithErrorJSON(c, http.StatusBadRequest, "Fail to take the order")
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "SUCCESS"})
}
