package usecase

import (
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/repository/apextransrepo"
)

type ApexTransUsecase interface {
	GetTabtransListApx(kuitansi string) ([]web.TabtransInfoApx, error)
	GetTabtransListByStanApx(stan string) ([]web.TabtransInfoApx, error)
}

type apextransUsecase struct{}

func NewApexTransUsecase() ApexTransUsecase {
	return &apextransUsecase{}
}

func (a *apextransUsecase) GetTabtransListApx(kuitansi string) (detailTx []web.TabtransInfoApx, er error) {
	repo, _ := apextransrepo.NewApexTransRepo()

	detailTx, er = repo.GetTabtransListApx(kuitansi)
	if er != nil {
		return detailTx, er
	}

	if len(detailTx) == 0 {
		return detailTx, err.NoRecord
	}

	return detailTx, nil
}

func (a *apextransUsecase) GetTabtransListByStanApx(stan string) (detailTx []web.TabtransInfoApx, er error) {
	repo, _ := apextransrepo.NewApexTransRepo()

	detailTx, er = repo.GetTabtransListByStanApx(stan)
	if er != nil {
		return detailTx, er
	}

	if len(detailTx) == 0 {
		return detailTx, err.NoRecord
	}

	return detailTx, nil
}
