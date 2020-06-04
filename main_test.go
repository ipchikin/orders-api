package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

// MainTestSuite is a testify suite
type MainTestSuite struct {
	suite.Suite
	// BM     models.BaseModel
	Router *gin.Engine
	W      *httptest.ResponseRecorder
}

// TestMainTestSuite is the core Test function
func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}

// SetupSuite will be executed before all tests are started
func (suite *MainTestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)
	suite.Router = setupRouter()
}

// BeforeTest will be executed before every test is started
func (suite *MainTestSuite) BeforeTest(suiteName, testName string) {
	// Reset the HTTP request
	suite.W = httptest.NewRecorder()

	// Make db connection with db choice
	// suite.BM.Connect(suite.dbType, suite.dbConString)
}

// AfterTest will be executed after every test is finished
func (suite *MainTestSuite) AfterTest(suiteName, testName string) {
}

// TearDownSuite will be executed after all tests are finished
func (suite *MainTestSuite) TearDownSuite() {
	// // Truncate orders table
	// result, err := suite.BM.DB.Exec("TRUNCATE TABLE orders")
	// suite.Nil(err)
	// suite.NotNil(result)
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
