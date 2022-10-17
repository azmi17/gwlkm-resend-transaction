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
	GetRetransTxInfo(stan string) (entities.RetransTxInfo, error)
	ChangeResponseCode(entities.ChangeResponseCode) error
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
	if data.ResponseCode != constant.Success {
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
	if reversedData.Ref_Stan == "" && reversedData.Response_Code != constant.Success {
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
		// Duplicate & assign new value from origin record..
		newTrx.Stan = "RT" + helper.GenerateSTAN()[2:12]
		newStan = newTrx.Stan
		newTrx.Ref_Stan = reversedData.Stan
		newTrx.Dc = "d"
		newTrx.Tgl_Trans_Str = reversedData.Tgl_Trans_Str
		newTrx.Ref = reversedData.Product_Code + newTrx.Stan
		newTrx.Response_Code = constant.Suspect

		// DO: DUPLICATE DATA
		er := dataRepo.DuplicatingData(newTrx)
		if er != nil {
			return newStan, er
		}

	}

	// Extracting TransHistory into MsgTransHistory..
	isoMsg := entities.MsgTransHistory{
		MTI:            newTrx.Mti,
		BankCode:       newTrx.Bank_Code,
		ProcessingCode: newTrx.Processing_Code,
		Stan:           newTrx.Stan,
		Date:           newTrx.Tgl_Trans_Str,
		Msg:            newTrx.Msg,
	}

	// RECYCLE REVERSED TRANSACTION
	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleReversedTransaction(&isoMsg)
	if isoMsg.ResponseCode == constant.Success {
		er = dataRepo.ChangeResponseCode(constant.Success, newTrx.Stan, 0)
		if er != nil {
			return newStan, err.InternalServiceError
		}
	} else {
		er = errors.New(isoMsg.Msg)
	}

	return newStan, er
}
func (r *retransactionUsecase) GetRetransTxInfo(stan string) (txInfo entities.RetransTxInfo, er error) {
	dataRepo, _ := datatransrepo.NewDatatransRepo()
	if txInfo, er = dataRepo.GetRetransTxInfo(stan); er != nil {
		return txInfo, er
	}
	return txInfo, nil
}

func (r *retransactionUsecase) ChangeResponseCode(payload entities.ChangeResponseCode) (er error) {
	dataRepo, _ := datatransrepo.NewDatatransRepo()
	er = dataRepo.ChangeResponseCode(payload.RC, payload.Stan, 0)
	if er != nil {
		return er
	}
	return nil
}
