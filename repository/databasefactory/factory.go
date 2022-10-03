package databasefactory

import (
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
	"os"
)

func GetDatabase() (db Database, err error) {

	driverName := os.Getenv("app.database_driver")

	// Default driver will return MYSQL-Driver
	if driverName == "" {
		driverName = drivers.MYSQL
	}

	if driverName == drivers.MYSQL {
		return newMysqlImpl(), nil
	} else {
		return db, errors.New("unimplement database driver")
	}
}
