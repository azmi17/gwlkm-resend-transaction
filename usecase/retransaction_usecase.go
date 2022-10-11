package usecase

import (
	"errors"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/constant"
	"gwlkm-resend-transaction/repository/datatransrepo"
	"gwlkm-resend-transaction/repository/retransactionepo"
)

type RetransactionUsecase interface {
	ResendTransaction(stan string) error
	ResendReversedTransaction(stan string) (string, error)
}

type retransactionUsecase struct{}

func NewRetransactionUsecase() RetransactionUsecase {
	return &retransactionUsecase{}
}

func (r *retransactionUsecase) ResendTransaction(stan string) (er error) {

	// INIT REPOSITORY BANK DATA
	dataRepo, _ := datatransrepo.NewDatatransRepo()

	// SEARCH DATA BY STAN
	var data entities.MsgTransHistory
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	// RECYCLE TRANSACTION
	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleTransaction(&data)
	if data.ResponseCode != "0000" {
		return errors.New(data.Msg)
	} else {
		return nil
	}
}

func (r *retransactionUsecase) ResendReversedTransaction(stan string) (newStan string, er error) {

	// REPEAT TRANS COND
	isRepeat := false

	// INIT REPOSITORY BANK DATA
	dataRepo, _ := datatransrepo.NewDatatransRepo()

	// SEARCH REVERSED DATA BY STAN
	var reversedData entities.TransHistory
	if reversedData, er = dataRepo.GetReversedData(stan); er != nil {
		return newStan, er
	}

	// Validation chk: If ref_stan == ""  & RC is not 0000
	if reversedData.Ref_Stan == "" && reversedData.Response_Code != "0000" {
		return newStan, err.RCMustBeSuccess
	}

	// Check if data repeated => isRepeat will be TRUE
	if reversedData.Ref_Stan != "" && (reversedData.Response_Code == constant.Suspect || reversedData.Response_Code == constant.Success) {
		isRepeat = true
	}
	newTrx := reversedData
	newStan = reversedData.Stan

	// If not repeat
	if !isRepeat {
		// Dupolicating & assign new value from origin record..
		newTrx.Stan = helper.GenerateSTAN()
		newStan = newTrx.Stan
		newTrx.Ref_Stan = reversedData.Stan
		newTrx.Tgl_Trans_Str = helper.GetCurrentDate()
		newTrx.Ref = reversedData.Product_Code + newTrx.Stan
		newTrx.Response_Code = constant.Suspect

		// DO: DUPLICATE DATA
		er := dataRepo.DuplicatingData(newTrx)
		if er != nil {
			return newStan, er
		}

		// CHANGE RC: ORIGIN RECORD
		er = dataRepo.ChangeRcOnReversedData(constant.Failed, reversedData.Stan)
		if er != nil {
			return newStan, err.InternalServiceError
		}
	}

	// Extracting TransHistory into MsgTransHistory..
	isoMsg := entities.MsgTransHistory{
		MTI:      newTrx.Mti,
		BankCode: newTrx.Bank_Code,
		Stan:     newTrx.Stan,
		Date:     newTrx.Tgl_Trans_Str,
		Msg:      newTrx.Msg,
	}

	// RECYCLE REVERSED TRANSACTION
	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleReversedTransaction(&isoMsg)
	if isoMsg.ResponseCode == constant.Success {
		er = dataRepo.ChangeRcOnReversedData(constant.Success, newTrx.Stan)
		if er != nil {
			return newStan, err.InternalServiceError
		}
	} else {
		er = errors.New(isoMsg.Msg)
	}

	return newStan, er
}
