package apextransrepo

import (
	"database/sql"
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
)

func NewApexTransRepo() (ApexTransRepo, error) {
	apexConn := databasefactory.Apex.GetConnection()
	currentDriver := databasefactory.Apex.GetDriverName()
	if currentDriver == drivers.MYSQL {
		return newApexTransRepoMysqlImpl(apexConn.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}
}
