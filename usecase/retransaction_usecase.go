package usecase

import (
	"errors"
	"gwlkm-resend-transaction/entities"
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

	var data entities.MsgTransHistory

	// INIT REPOSITORY BANK DATA
	dataRepo, _ := datatransrepo.NewDatatransRepo()

	// CALL DATA BY STAN
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	// CALL RECYCLE TRANSACTION
	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleTransaction(&data)
	if data.ResponseCode != "0000" {
		return errors.New(data.Msg)
	} else {
		return nil
	}
}

func (r *retransactionUsecase) ResendReversedTransaction(stan string) (refStan string, er error) {

	var reversedData entities.TransHistory

	// INIT REPOSITORY BANK DATA
	dataRepo, _ := datatransrepo.NewDatatransRepo()

	// CALL REVERSED DATA BY STAN
	if reversedData, er = dataRepo.GetReversedData(stan); er != nil {
		return refStan, er
	}

	// RECOMPOSING FOR DUPLICATE PROCS
	newTrx := entities.TransHistory{}
	newTrx.Stan = helper.GenerateSTAN()
	newTrx.Ref_Stan = reversedData.Stan
	newTrx.Tgl_Trans_Str = helper.GetCurrentDate()
	// newTrx.Tgl_Trans_Str = reversedData.Tgl_Trans_Str
	newTrx.Bank_Code = reversedData.Bank_Code
	newTrx.Rek_Id = reversedData.Rek_Id
	newTrx.Mti = reversedData.Mti
	newTrx.Processing_Code = reversedData.Processing_Code
	newTrx.Biller_Code = reversedData.Biller_Code
	newTrx.Product_Code = reversedData.Product_Code
	newTrx.Subscriber_Id = reversedData.Subscriber_Id
	newTrx.Dc = reversedData.Dc
	newTrx.Response_Code = "0000"
	newTrx.Amount = reversedData.Amount
	newTrx.Qty = reversedData.Qty
	newTrx.Profit_Included = reversedData.Profit_Included
	newTrx.Profit_Excluded = reversedData.Profit_Excluded
	newTrx.Profit_Share_Biller = reversedData.Profit_Share_Biller
	newTrx.Profit_Share_Aggregator = reversedData.Profit_Share_Aggregator
	newTrx.Profit_Share_Bank = reversedData.Profit_Share_Bank
	newTrx.Markup_Total = reversedData.Markup_Total
	newTrx.Markup_Share_Aggregator = reversedData.Markup_Share_Aggregator
	newTrx.Markup_Share_Bank = reversedData.Markup_Share_Bank
	newTrx.Msg = reversedData.Msg
	newTrx.Msg_Response = reversedData.Msg_Response
	newTrx.Bit39_Bit48_Hulu = reversedData.Bit39_Bit48_Hulu
	newTrx.Saldo_Before_Trans = reversedData.Saldo_Before_Trans
	newTrx.Keterangan = reversedData.Keterangan
	newTrx.Ref = reversedData.Product_Code + newTrx.Stan
	newTrx.Synced_Ibs_Core = reversedData.Synced_Ibs_Core
	newTrx.Synced_Ibs_Core_Description = reversedData.Synced_Ibs_Core_Description
	newTrx.Bris_Original_Data = reversedData.Bris_Original_Data
	newTrx.Gateway_Id = reversedData.Gateway_Id
	newTrx.Id_User = reversedData.Id_User
	newTrx.Id_Raw = reversedData.Id_Raw
	newTrx.Advice_Count = reversedData.Advice_Count
	newTrx.Status_Id = reversedData.Status_Id
	newTrx.Nohp_Notif = reversedData.Nohp_Notif
	newTrx.Score = reversedData.Score
	newTrx.No_Hp_Alternatif = reversedData.No_Hp_Alternatif
	newTrx.Inc_Notif_Status = reversedData.Inc_Notif_Status
	newTrx.Fee_Rek_Induk = reversedData.Fee_Rek_Induk

	// CALL DUPLICATE DATA
	cpy, er := dataRepo.DuplicatingData(newTrx)
	if er != nil {
		return refStan, er
	}

	// if newTrx.Ref_Stan == reversedData.Stan {
	// 	return refStan, err.DuplicateEntry
	// }

	// Extracting TransHistory into MsgTransHistory..
	isoMsg := entities.MsgTransHistory{
		MTI:      cpy.Mti,
		BankCode: cpy.Bank_Code,
		Stan:     cpy.Stan,
		Ref:      cpy.Ref_Stan,
		Date:     cpy.Tgl_Trans_Str,
		Msg:      cpy.Msg,
	}

	// CALL RECYCLE REVERSED TRANSACTION
	reTransRepo := retransactionepo.NewRetransactionRepo()
	reTransRepo.RecycleReversedTransaction(&isoMsg)
	if isoMsg.ResponseCode != "0000" {
		er = dataRepo.ChangeRcOnReversedData(constant.Pending, cpy.Stan)
		if er != nil {
			return refStan, er
		}
		return refStan, errors.New(isoMsg.Msg)
	} else {
		refStan = newTrx.Product_Code + newTrx.Stan
		return refStan, nil
	}
}
