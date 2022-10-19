package datatransrepo

import (
	"gwlkm-resend-transaction/entities"
)

type DatatransRepo interface {
	// Echannel
	GetData(stan string) (entities.MsgTransHistory, error)
	GetServeAddr(bankCode string) (entities.CoreAddr, error)
	GetReversedData(stan string) (entities.TransHistory, error)
	DuplicatingData(copy entities.TransHistory) error
	ChangeResponseCode(rc, stan string, transId int) error
	AddStanReference(copy entities.StanReference) error
	GetRetransTxInfo(stan string) (entities.RetransTxInfo, error)

	// Apex
	GetTransIdApx() (int, error)
	GetTxInfoApx(kuitansi string) (entities.TransApx, error)
	DuplicatingTxApx(copy entities.TransApx) error
	DeleteTxApx(kuitansi string) error

	RecycleTxApx(kuitansi string) error
}
