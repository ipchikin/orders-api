package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// BaseModel of other models
type BaseModel struct {
	DB *sqlx.DB
}

// Connect model to db
func (bm *BaseModel) Connect(driver, user, password, host, port, database string, maxIdleConns int) (err error) {
	bm.DB, err = sqlx.Connect(driver,
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			user,
			password,
			host,
			port,
			database,
		),
	)
	if err != nil {
		return
	}

	bm.DB.SetMaxIdleConns(maxIdleConns)
	return
}
