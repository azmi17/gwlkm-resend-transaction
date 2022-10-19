package datatransrepo

import (
	"database/sql"
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
)

func NewDatatransRepo() (DatatransRepo, error) {

	conn1 := databasefactory.AppDb1.GetConnection()
	currentDriver := databasefactory.AppDb1.GetDriverName()
	if currentDriver == drivers.MYSQL {
		conn2 := databasefactory.AppDb2.GetConnection()
		return newDatatransRepoMysqlImpl(conn1.(*sql.DB), conn2.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}

}
