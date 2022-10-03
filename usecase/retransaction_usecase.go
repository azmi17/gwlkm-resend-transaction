package usecase

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/repository/datatransrepo"
	"gwlkm-resend-transaction/repository/retransactionepo"
)

type RetransactionUsecase interface {
	ResendTransaction(stan string) error
}

type retransactionUsecase struct{}

func NewRetransactionUsecase() RetransactionUsecase {
	return &retransactionUsecase{}
}

func (e *retransactionUsecase) ResendTransaction(stan string) (er error) {

	var data entities.TransHisotry

	dataRepo, _ := datatransrepo.NewDatatransRepo()
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleTransaction(data)
	// return reTransRepo.RecycleTransaction(data)
	return nil
}
