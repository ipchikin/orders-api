package models

import (
	"orders-api/configs"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
)

// OrdersTestSuite
type OrdersTestSuite struct {
	suite.Suite
	OM            OrdersModel
	SampleOrderID string
}

// TestOrdersTestSuite runs all tests
func TestOrdersTestSuite(t *testing.T) {
	suite.Run(t, new(OrdersTestSuite))
}

// SetupSuite
func (suite *OrdersTestSuite) SetupSuite() {
	// Load .env
	_, caller, _, ok := runtime.Caller(0)
	suite.True(ok)
	suite.Nil(godotenv.Load(filepath.Join(filepath.Dir(caller), "..") + "/.env"))

	// Load test config
	gin.SetMode(gin.TestMode)
	cfg, err := configs.LoadConfig("test")
	suite.Nil(err)

	// Connect to test DB
	err = suite.OM.Connect(
		cfg.DBConfig.Driver,
		cfg.DBConfig.User,
		cfg.DBConfig.Password,
		cfg.DBConfig.Host,
		cfg.DBConfig.Port,
		cfg.DBConfig.Database,
		cfg.DBConfig.MaxIdleConns,
	)
	suite.Nil(err)

	suite.SampleOrderID = "15e899be-6b67-47e3-8290-857524696184"
}

// BeforeTest will be executed before every test
func (suite *OrdersTestSuite) BeforeTest(suiteName, testName string) {
	// Truncate orders table
	_, err := suite.OM.DB.Exec("TRUNCATE TABLE orders")
	suite.Nil(err)

	// Create sample order
	res, err := suite.OM.DB.Exec(`INSERT INTO orders (id, origin_lat, origin_long, destination_lat, destination_long, distance, status) VALUES (?, '22.281980', '114.161370', '22.318359', '114.157913', 7635, 'UNASSIGNED')`, suite.SampleOrderID)
	suite.Nil(err)
	affected, err := res.RowsAffected()
	suite.Nil(err)
	suite.Equal(int64(1), affected)
}

// TearDownSuite will be executed after all tests
func (suite *OrdersTestSuite) TearDownSuite() {
	// Truncate orders table
	_, err := suite.OM.DB.Exec("TRUNCATE TABLE orders")
	suite.Nil(err)
}

// TestPlace
func (suite *OrdersTestSuite) TestPlace() {
	id, err := suite.OM.Place([2]string{"22.281980", "114.161370"}, [2]string{"22.318359", "114.157913"}, 7635, "UNASSIGNED")
	suite.Nil(err)
	suite.NotNil(id)
}

// TestTake
func (suite *OrdersTestSuite) TestTake() {
	suite.Nil(suite.OM.Take(suite.SampleOrderID, "TAKEN"))
}

// TestTakeNotExists with not exists order id
func (suite *OrdersTestSuite) TestTakeNotExists() {
	suite.NotNil(suite.OM.Take("00000000-0000-0000-0000-000000000000", "TAKEN"))
}

// TestTakeNotUnassigned with not unassigned order status
func (suite *OrdersTestSuite) TestTakeNotUnassigned() {
	_, err := suite.OM.DB.Exec("UPDATE orders SET status = ? WHERE id = ?", "TAKEN", suite.SampleOrderID)
	suite.Nil(err)
	suite.NotNil(suite.OM.Take(suite.SampleOrderID, "TAKEN"))
}

// TestTakeRace with race condition
func (suite *OrdersTestSuite) TestTakeRace() {
	tx, err := suite.OM.DB.Beginx()
	suite.Nil(err)

	// Lock the order
	var order Order
	suite.Nil(tx.Get(&order, "SELECT id, status FROM orders WHERE id=? FOR UPDATE", suite.SampleOrderID))

	// Take the locked order
	suite.NotNil(suite.OM.Take(suite.SampleOrderID, "TAKEN"))

	suite.Nil(tx.Rollback())
}

// TestList
func (suite *OrdersTestSuite) TestList() {
	orders, err := suite.OM.List(1, 10)
	suite.Nil(err)
	suite.NotNil(orders)
}
