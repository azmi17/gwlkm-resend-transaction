package apextransrepo

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
)

type ApexTransRepo interface {
	GetTransIdApx() (int, error)
	GetTabtransTxInfoApx(kuitansi, bankCode string) (entities.TransApx, error)
	DuplicatingTxApx(copy ...entities.TransApx) error
	DuplicatingUnitTestTxApx(copy ...entities.TransApx) error
	DeleteTxApx(kuitansi, bankCode string) error
	GetTabtransListApx(kuitansi string) ([]web.TabtransInfoApx, error)
	GetCreditTransferSMLkmApx(kuitansi, MyKdTrans, bankCode string) (entities.TransApx, error)
	DuplicateCreditTransferSMLkmApx(copy entities.TransApx) error // DuplicateCreditTransferSMLkmApx(copy ...entities.TransApx) error
	GetTabtransListByStanApx(stan string) ([]web.TabtransInfoApx, error)
}
