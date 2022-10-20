package echanneltransrepo

import (
	"database/sql"
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
)

func NewEchannelTransRepo() (EchannelTransRepo, error) {
	echannelConn := databasefactory.Echannel.GetConnection()
	currentDriver := databasefactory.Echannel.GetDriverName()
	if currentDriver == drivers.MYSQL {
		return newEchannelTransRepoMysqlImpl(echannelConn.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}
}
