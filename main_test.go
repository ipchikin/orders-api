package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"orders-api/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// MainTestSuite
type MainTestSuite struct {
	suite.Suite
	BM            models.BaseModel
	Router        *gin.Engine
	W             *httptest.ResponseRecorder
	SampleOrderID string
}

// TestMainTestSuite runs all tests
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

// SetupSuite
func (suite *MainTestSuite) SetupSuite() {
	// Load test config
	gin.SetMode(gin.TestMode)
	cfg, err := loadConfig()
	suite.Nil(err)

	// Connect to test DB
	err = suite.BM.Connect(
		cfg.DBConfig.Driver,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Database,
		cfg.DBConfig.MaxIdleConns,
	)
	suite.Nil(err)

	suite.Router = setupRouter()

	suite.SampleOrderID = "15e899be-6b67-47e3-8290-857524696184"
}

// BeforeTest will be executed before every test
func (suite *MainTestSuite) BeforeTest(suiteName, testName string) {
	// Truncate orders table
	_, err := suite.BM.DB.Exec("TRUNCATE TABLE orders")
	suite.Nil(err)

	// Create sample order
	res, err := suite.BM.DB.Exec(`INSERT INTO orders (id, origin_lat, origin_long, destination_lat, destination_long, distance, status) VALUES (?, '22.281980', '114.161370', '22.318359', '114.157913', 7635, 'UNASSIGNED')`, suite.SampleOrderID)
	suite.Nil(err)
	affected, err := res.RowsAffected()
	suite.Nil(err)
	suite.Equal(int64(1), affected)

	// Reset the HTTP request
	suite.W = httptest.NewRecorder()
}

// TestPlaceOrder with correct coordinates
func (suite *MainTestSuite) TestPlaceOrder() {
	var jsonStr = []byte(`{"origin":["22.281980","114.161370"],"destination":["22.318359","114.157913"]}`)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	suite.Router.ServeHTTP(suite.W, req)

	suite.Equal(http.StatusOK, suite.W.Code)
}

// TestPlaceOrderOverDecimalPlaces with correct coordinates but over 6 decimal places
func (suite *MainTestSuite) TestPlaceOrderOverDecimalPlaces() {
	var jsonStr = []byte(`{"origin":["22.281980123456","114.161370123456"],"destination":["22.318359123456","114.157913123456"]}`)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	suite.Router.ServeHTTP(suite.W, req)

	suite.Equal(http.StatusOK, suite.W.Code)
}

// TestTakeOrder with correct id
func (suite *MainTestSuite) TestTakeOrder() {
	var jsonStr = []byte(`{"status": "TAKEN"}`)
	req, _ := http.NewRequest("PATCH", "/orders/"+suite.SampleOrderID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	suite.Router.ServeHTTP(suite.W, req)

	suite.Equal(http.StatusOK, suite.W.Code)
}

// TestTakeOrder with race condition
func (suite *MainTestSuite) TestTakeOrderRace() {
	tx, err := suite.BM.DB.Beginx()
	suite.Nil(err)

	// Lock the order
	var order models.Order
	suite.Nil(tx.Get(&order, "SELECT id, status FROM orders WHERE id=? FOR UPDATE", suite.SampleOrderID))

	var jsonStr = []byte(`{"status": "TAKEN"}`)
	req, _ := http.NewRequest("PATCH", "/orders/"+suite.SampleOrderID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	suite.Router.ServeHTTP(suite.W, req)

	suite.Equal(http.StatusBadRequest, suite.W.Code)

	suite.Nil(tx.Rollback())
}

// TestListOrders with correct page and limit
func (suite *MainTestSuite) TestListOrders() {
	req, _ := http.NewRequest("GET", "/orders?page=1&limit=10", nil)
	req.Header.Set("Content-Type", "application/json")
	suite.Router.ServeHTTP(suite.W, req)

	suite.Equal(http.StatusOK, suite.W.Code)
}
