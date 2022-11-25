package sysusertempsrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
)

func newSysUserMysqlImpl(apexConn *sql.DB) SysUserRepo {
	return &SysUserMysqlImpl{
		apexDb: apexConn,
	}
}

type SysUserMysqlImpl struct {
	apexDb *sql.DB
}

func (s *SysUserMysqlImpl) ResetUserPassword(user entities.SysDaftarUser) (sysUser entities.SysDaftarUser, er error) {

	thisRepo, _ := NewSysUserRepo()
	_, er = thisRepo.FindByUserName(user.User_Name)
	if er != nil {
		return sysUser, err.NoRecord
	}

	stmt, er := s.apexDb.Prepare("UPDATE sys_daftar_user SET user_web_password = ? WHERE user_name = ?")
	if er != nil {
		return sysUser, errors.New(fmt.Sprint("error while prepare update apex user web password: ", er.Error()))
	}

	defer func() {
		_ = stmt.Close()
	}()

	if _, er := stmt.Exec(user.User_Web_Password, user.User_Name); er != nil {
		return sysUser, errors.New(fmt.Sprint("error while update apex user web password: ", er.Error()))
	}

	return user, nil
}

func (s *SysUserMysqlImpl) FindByUserName(userName string) (user entities.SysDaftarUser, er error) {
	row := s.apexDb.QueryRow(`SELECT
		user_id, 
		user_name,
		nama_lengkap
	FROM sys_daftar_user WHERE user_name = ? LIMIT 1`, userName)
	er = row.Scan(
		&user.User_Id,
		&user.User_Name,
		&user.Nama_Lengkap,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return user, err.NoRecord
		} else {
			return user, errors.New(fmt.Sprint("error while get user name: ", er.Error()))
		}
	}
	return
}
