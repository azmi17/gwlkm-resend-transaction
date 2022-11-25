package usecase

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/sysusertempsrepo"
)

type SysUserUsecase interface {
	ResetSysUserPassword(KodeLkm entities.KodeLKMFilter) (entities.ResetApexPwdResponse, error)
}

type sysUserUsecase struct{}

func NewSysUserUsecase() SysUserUsecase {
	return &sysUserUsecase{}
}

func (s *sysUserUsecase) ResetSysUserPassword(KodeLkm entities.KodeLKMFilter) (resp entities.ResetApexPwdResponse, er error) {

	if KodeLkm.KodeLkm == "" {
		return resp, err.FieldMustBeExist
	}

	sysUserRepo, _ := sysusertempsrepo.NewSysUserRepo()

	sysDaftarUser := entities.SysDaftarUser{}
	sysDaftarUser.User_Name = KodeLkm.KodeLkm
	sysDaftarUser.User_Web_Password_Hash, sysDaftarUser.User_Web_Password = helper.HashSha1Pass()

	if sysDaftarUser, er = sysUserRepo.ResetUserPassword(sysDaftarUser); er != nil {
		return resp, er
	}

	updResp := entities.ResetApexPwdResponse{}
	updResp.KodeLkm = sysDaftarUser.User_Name
	updResp.Password_Smec = sysDaftarUser.User_Web_Password_Hash

	return updResp, nil
}
