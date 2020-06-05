package orders

import (
	"net/http"
	"orders-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListOrders controller
func ListOrders(c *gin.Context) {
	q := c.Request.URL.Query()

	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid page value")
		return
	}

	if page <= 0 {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid page value")
		return
	}

	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid limit value")
		return
	}

	if limit < 0 {
		abortWithErrorJSON(c, http.StatusBadRequest, "Invalid limit value")
		return
	}

	// List orders
	om := models.OrdersModel{BaseModel: c.MustGet("db").(models.BaseModel)}
	orders, err := om.List(page, limit)
	if err != nil {
		abortWithErrorJSON(c, http.StatusBadRequest, "Fail to list the orders")
		return
	}

	c.JSON(http.StatusOK, orders)
}
