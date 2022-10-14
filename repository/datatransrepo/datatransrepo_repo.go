package datatransrepo

import (
	"gwlkm-resend-transaction/entities"
)

type DatatransRepo interface {
	GetData(stan string) (entities.MsgTransHistory, error)
	GetServeAddr(bankCode string) (entities.CoreAddr, error)
	GetReversedData(stan string) (entities.TransHistory, error)
	DuplicatingData(copy entities.TransHistory) error
	ChangeResponseCode(rc, stan string, transId int) error
	AddStanReference(copy entities.StanReference) error
	GetRetransTxInfo(stan string) (entities.RetransTxInfo, error)
}
