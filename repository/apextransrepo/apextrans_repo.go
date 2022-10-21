package apextransrepo

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
)

type ApexTransRepo interface {
	GetTransIdApx() (int, error)
	GetTxInfoApx(kuitansi string) (entities.TransApx, error)
	DuplicatingTxApx(copy ...entities.TransApx) error
	DuplicatingUnitTestTxApx(copy ...entities.TransApx) error
	DeleteTxApx(kuitansi string) error
	GetTabtransApx(kuitansi string) ([]web.TabtransInfoApx, error)
}
