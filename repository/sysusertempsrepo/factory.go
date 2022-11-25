package sysusertempsrepo

import (
	"database/sql"
	"errors"
	"gwlkm-resend-transaction/repository/databasefactory"
	"gwlkm-resend-transaction/repository/databasefactory/drivers"
)

func NewSysUserRepo() (SysUserRepo, error) {
	apexConn := databasefactory.SysApex.GetConnection()
	currentDriver := databasefactory.SysApex.GetDriverName()
	if currentDriver == drivers.MYSQL {
		return newSysUserMysqlImpl(apexConn.(*sql.DB)), nil
	} else {
		return nil, errors.New("unimplemented database driver")
	}
}
