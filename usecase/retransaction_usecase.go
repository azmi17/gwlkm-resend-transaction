package usecase

import (
	"errors"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/constant"
	"gwlkm-resend-transaction/repository/echanneltransrepo"
	"gwlkm-resend-transaction/repository/retransactionepo"
)

type RetransactionUsecase interface {
	ResendTransaction(stan string) error
	ResendGwlkmTransaction(stan string) (string, error)
	GetRetransTxInfo(stan string) (web.RetransTxInfo, error)
	ResendLkmTransferSMprematureRevOnCre(string) (er error)
	ResendReversalBeforeRecycleGwlkmTrx(string) (er error)
}

type retransactionUsecase struct{}

func NewRetransactionUsecase() RetransactionUsecase {
	return &retransactionUsecase{}
}

func (r *retransactionUsecase) ResendTransaction(stan string) (er error) {

	// INIT REPOSITORY BANK DATA
	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()

	// SEARCH DATA BY STAN
	var data entities.IsoMessageBody
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

func (r *retransactionUsecase) ResendGwlkmTransaction(stan string) (newStan string, er error) {
	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	reTransRepo := retransactionepo.NewRetransactionRepo()

	// IS REPEAT CONDITION
	isRepeat := false

	// Search Origin Data By STAN
	var originData entities.TransHistory
	if originData, er = dataRepo.GetOriginData(stan); er != nil {
		return newStan, er
	}

	// Validation chk: If ref_stan == ""  & RC is not 0000
	if originData.Ref_Stan == "" && originData.Response_Code != constant.Success {
		return newStan, err.RCMustBeSuccess
	}

	// Check if data repeated => isRepeat will be true
	if originData.Ref_Stan != "" && (originData.Response_Code == constant.Suspect || originData.Response_Code == constant.Success) {
		isRepeat = true
	}
	newTrx := originData
	newStan = originData.Stan

	// If the condition is false, the statement below will be skipped
	oldStan := originData.Stan
	if !isRepeat {
		// assign new value from origin record..
		newTrx.Stan = "RT" + helper.GenerateSTAN()[2:12]
		newStan = newTrx.Stan
		newTrx.Ref_Stan = originData.Stan
		newTrx.Dc = "d"
		newTrx.Tgl_Trans_Str = originData.Tgl_Trans_Str
		newTrx.Ref = originData.Product_Code + newTrx.Stan
		newTrx.Response_Code = constant.Suspect

		// DO: Duplicate origin data
		er = dataRepo.DuplicatingData(newTrx)
		if er != nil {
			return newStan, er
		}
	} else {
		oldStan = newTrx.Ref_Stan
	}

	// Extracting TransHistory into IsoMessageBody..
	isoMsg := entities.IsoMessageBody{
		MTI:            newTrx.Mti,
		BankCode:       newTrx.Bank_Code,
		ProcessingCode: newTrx.Processing_Code,
		Stan:           newTrx.Stan,
		Date:           newTrx.Tgl_Trans_Str,
		Msg:            newTrx.Msg,
	}

	// DO: Send reversal before RecycleGwlkmTransaction()
	er = reTransRepo.ResendReversalBeforeRecycleGwlkmTransaction(oldStan)
	if er != nil {
		entities.PrintError(er.Error())
	}

	// DO: Recycle Transaction
	reTransRepo.RecycleGwlkmTransaction(&isoMsg)
	if isoMsg.ResponseCode == constant.Success {
		er = dataRepo.ChangeResponseCode(constant.Success, newTrx.Stan, 0)
		if er != nil {
			return newStan, err.InternalServiceError
		}
	} else {
		if isoMsg.Msg == helper.AlreadyTransacted {
			er = dataRepo.ChangeResponseCode(constant.Success, newTrx.Stan, 0)
			if er != nil {
				return newStan, err.InternalServiceError
			}
		}
		er = errors.New(isoMsg.Msg)
	}
	return newStan, er
}

func (r *retransactionUsecase) GetRetransTxInfo(stan string) (txInfo web.RetransTxInfo, er error) {
	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	if txInfo, er = dataRepo.GetRetransTxInfo(stan); er != nil {
		return txInfo, er
	}
	return txInfo, nil
}

func (r *retransactionUsecase) ResendLkmTransferSMprematureRevOnCre(stan string) (er error) {

	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	var data entities.IsoMessageBody
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleLkmTransferSMprematureRevOnCre(&data)
	if data.ResponseCode != constant.Success {
		return errors.New(data.Msg)
	} else {
		return nil
	}
}

func (r *retransactionUsecase) ResendReversalBeforeRecycleGwlkmTrx(stan string) (er error) {

	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	var data entities.IsoMessageBody
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.ResendReversalGwlkmTransaction(&data)
	if data.ResponseCode != constant.Success {
		return errors.New(data.Msg)
	} else {
		return nil
	}
}
