package datatransrepo

import (
	"database/sql"
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
)

func NewDatatransRepo() (DatatransRepo, error) {

	conn := databasefactory.AppDb.GetConnection()

	currentDriver := databasefactory.AppDb.GetDriverName()
	if currentDriver == drivers.MYSQL {
		return newDatatransRepoMysqlImpl(conn.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}

}
