package sysusertempsrepo

import "gwlkm-resend-transaction/entities"

type SysUserRepo interface {
	FindByUserName(userName string) (entities.SysDaftarUser, error)
	ResetUserPassword(user entities.SysDaftarUser) (entities.SysDaftarUser, error)
}
