package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const driverName = "mysql"
const dataSourceName = "root:secret@tcp(localhost:3306)/orders"
const maxIdleConns = 2

// TestConnect tests connecting to database
func TestConnect(t *testing.T) {
	bm := new(BaseModel)
	err := bm.Connect(driverName, "", maxIdleConns)
	assert.NotNil(t, err)

	err = bm.Connect(driverName, dataSourceName, maxIdleConns)
	assert.Nil(t, err)
}
