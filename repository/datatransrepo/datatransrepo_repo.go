package datatransrepo

import (
	"gwlkm-resend-transaction/entities"
)

type DatatransRepo interface {
	GetData(stan string) (entities.TransHisotry, error)
	GetServeAddr(bankCode string) (entities.CoreAddr, error)
}
