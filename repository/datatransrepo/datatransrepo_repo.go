package datatransrepo

import (
	"gwlkm-resend-transaction/entities"
)

type DatatransRepo interface {
	GetData(stan string) (entities.MsgTransHistory, error)
	GetServeAddr(bankCode string) (entities.CoreAddr, error)

	GetReversedData(stan string) (entities.TransHistory, error)
	DuplicatingData(duplicated entities.TransHistory) (entities.TransHistory, error)
	ChangeRcOnReversedData(stan string) error
}
