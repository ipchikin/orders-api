package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// BaseModel
type BaseModel struct {
	DB *sqlx.DB
}

// Connect to db
func (bm *BaseModel) Connect(driverName, dataSourceName string, maxIdleConns int) (err error) {
	bm.DB, err = sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return
	}

	bm.DB.SetMaxIdleConns(maxIdleConns)
	return
}
