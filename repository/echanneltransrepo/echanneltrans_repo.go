package echanneltransrepo

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
)

type EchannelTransRepo interface {
	GetData(stan string) (entities.IsoMessageBody, error)
	GetServeAddr(bankCode string) (entities.CoreAddrInfo, error)
	GetOriginData(stan string) (entities.TransHistory, error)
	DuplicatingData(copy entities.TransHistory) error
	ChangeResponseCode(rc, stan string, transId int) error
	AddStanReference(copy entities.StanReferences) error
	GetRetransTxInfo(stan string) (web.RetransTxInfo, error)
	UpdateIsoMsg(msg, stan string) error
}
