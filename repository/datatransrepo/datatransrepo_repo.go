package datatransrepo

import (
	"gwlkm-resend-transaction/entities"
)

type DatatransRepo interface {
	GetData(stan string) (entities.MsgTransHistory, error)
	GetServeAddr(bankCode string) (entities.CoreAddr, error)
	GetReversedData(stan string) (entities.TransHistory, error)
	DuplicatingData(copy entities.TransHistory) error
	ChangeRcOnReversedData(rc, stan string) error
	RollbackDuplicate(stan string) error
}
