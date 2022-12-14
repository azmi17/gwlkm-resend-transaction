package retransactionepo

import (
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/apextransrepo"
	"gwlkm-resend-transaction/repository/constant"
	"gwlkm-resend-transaction/repository/echanneltransrepo"

	iso8583uParser "github.com/randyardiansyah25/iso8583u/parser"
	"github.com/randyardiansyah25/libpkg/net/tcp"
)

func NewRetransactionRepo() RetransactionRepo {
	return &retransactionRepoImpl{}
}

type retransactionRepoImpl struct {
}

func (r *retransactionRepoImpl) RecycleTransaction(dataTrans *entities.IsoMessageBody) (err error) {

	//TODO: Send ISO data to IP & Port GWLKM

	// ISO OBJ INIT
	isoUnMarshal, err := iso8583uParser.NewISO8583U()
	if err != nil {
		entities.PrintError("load package error", err.Error())
		return
	}
	// UNMARSHAL | RE-COMPOSE ISO
	err = isoUnMarshal.GoUnMarshal(dataTrans.Msg)
	if err != nil {
		entities.PrintError(err.Error())
		return
	}
	isoUnMarshal.SetMti(dataTrans.MTI)
	isoUnMarshal.SetField(3, "400700") // => 400700 extract current date (?)
	isoUnMarshal.SetField(4, isoUnMarshal.GetField(4))
	isoUnMarshal.SetField(5, isoUnMarshal.GetField(5))
	isoUnMarshal.SetField(6, isoUnMarshal.GetField(6))
	isoUnMarshal.SetField(7, isoUnMarshal.GetField(7))
	isoUnMarshal.SetField(8, isoUnMarshal.GetField(8))
	isoUnMarshal.SetField(11, isoUnMarshal.GetField(11))
	isoUnMarshal.SetField(12, isoUnMarshal.GetField(12))
	isoUnMarshal.SetField(13, isoUnMarshal.GetField(13))
	isoUnMarshal.SetField(18, isoUnMarshal.GetField(18))
	isoUnMarshal.SetField(26, isoUnMarshal.GetField(26))
	isoUnMarshal.SetField(32, isoUnMarshal.GetField(32))
	isoUnMarshal.SetField(37, isoUnMarshal.GetField(37))
	isoUnMarshal.SetField(40, isoUnMarshal.GetField(40))
	isoUnMarshal.SetField(41, isoUnMarshal.GetField(41))
	isoUnMarshal.SetField(42, isoUnMarshal.GetField(42))
	isoUnMarshal.SetField(43, isoUnMarshal.GetField(43))
	isoUnMarshal.SetField(47, isoUnMarshal.GetField(47))
	isoUnMarshal.SetField(61, isoUnMarshal.GetField(61))
	isoUnMarshal.SetField(100, isoUnMarshal.GetField(100))
	isoUnMarshal.SetField(103, isoUnMarshal.GetField(103))
	isoUnMarshal.SetField(104, "TEXTDB")

	// MARSHAL PROCS
	isoMsg, err := isoUnMarshal.GoMarshal()
	if err != nil {
		entities.PrintError(err.Error())
		return
	}

	// CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, err := repo.GetServeAddr(dataTrans.BankCode)
	if err != nil {
		entities.PrintError(err.Error())
	}

	// INIT TCP OBJ
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	entities.PrintLog("SEND:\n", isoUnMarshal.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			entities.PrintError(err.Error())
			return
		}
		// Override below:
		dataTrans.ResponseCode = isoUnMarshal.GetField(39)
		dataTrans.Msg = isoUnMarshal.GetField(48)
		entities.PrintLog("RECV:\n", isoUnMarshal.PrettyPrint())
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaction failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}
	return nil
}

func (r *retransactionRepoImpl) RecycleGwlkmTransaction(dataTrans *entities.IsoMessageBody) (err error) {

	//TODO: Send ISO data to IP & Port GWLKM

	// ISO OBJ
	iso, err := iso8583uParser.NewISO8583U()
	if err != nil {
		entities.PrintError("load package error", err.Error())
		return
	}

	// UNMARSHAL | RE-COMPOSE ISO
	err = iso.GoUnMarshal(dataTrans.Msg)
	if err != nil {
		entities.PrintError(err.Error())
		return
	}
	/*
		iso.SetMti(dataTrans.MTI)
		iso.SetField(3, dataTrans.ProcessingCode)
		iso.SetField(4, iso.GetField(4))
		iso.SetField(5, iso.GetField(5))
		iso.SetField(6, iso.GetField(6))
		iso.SetField(7, iso.GetField(7))
		iso.SetField(8, iso.GetField(8))
		iso.SetField(11, dataTrans.Stan)
		iso.SetField(12, helper.GetCurrentDate())
		iso.SetField(13, helper.GETMMDD[4:8])
		iso.SetField(18, iso.GetField(18))
		iso.SetField(26, iso.GetField(26))
		iso.SetField(32, iso.GetField(32))
		iso.SetField(37, iso.GetField(37))
		iso.SetField(40, iso.GetField(40))
		iso.SetField(41, iso.GetField(41))
		iso.SetField(42, iso.GetField(42))
		iso.SetField(43, iso.GetField(43))
		iso.SetField(47, iso.GetField(47))
		iso.SetField(61, iso.GetField(61))
		iso.SetField(100, iso.GetField(100))
		iso.SetField(103, iso.GetField(103))
		iso.SetField(104, iso.GetField(104))
	*/
	iso.SetMti(dataTrans.MTI)
	iso.SetField(3, dataTrans.ProcessingCode)
	iso.SetField(11, dataTrans.Stan)
	iso.SetField(12, helper.GetCurrentDate())
	iso.SetField(13, iso.GetField(13))

	// MARSHAL PROCS
	isoMsg, err := iso.GoMarshal()
	if err != nil {
		entities.PrintError(err.Error())
		return
	}

	// GET CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, err := repo.GetServeAddr(dataTrans.BankCode)
	if err != nil {
		entities.PrintError(err.Error())
	}

	// TCP OBJ INIT
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 80)
	entities.PrintLog("SEND:\n", iso.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL ISO FROM SENDER
	if st.Code == tcp.CONNOK {
		err = iso.GoUnMarshal(st.Message)
		if err != nil {
			entities.PrintError(err.Error())
			return
		}
		// Override below:
		dataTrans.ResponseCode = iso.GetField(39)
		dataTrans.Msg = iso.GetField(48)
		entities.PrintLog("RECV:\n", iso.PrettyPrint())
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaction failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	return nil
}

func (r *retransactionRepoImpl) RecycleLkmTransferSMprematureRevOnCre(dataTrans *entities.IsoMessageBody) (er error) {

	//TODO: Send ISO data to IP & Port GWLKM
	// ISO OBJ INIT
	iso, er := iso8583uParser.NewISO8583U()
	if er != nil {
		entities.PrintError("load package error", er.Error())
		return
	}

	// UNMARSHAL | RE-COMPOSE ISO
	er = iso.GoUnMarshal(dataTrans.Msg)
	if er != nil {
		entities.PrintError(er.Error())
		return
	}
	iso.SetMti(dataTrans.MTI)
	iso.SetField(3, "400700")
	iso.SetField(12, helper.GetCurrentDate())
	iso.SetField(104, "TINTCR")

	// MARSHAL PROCS
	isoMsg, er := iso.GoMarshal()
	if er != nil {
		entities.PrintError(er.Error())
		return
	}

	// CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, er := repo.GetServeAddr(dataTrans.BankCode)
	if er != nil {
		entities.PrintError(er.Error())
	}

	// INIT TCP OBJ
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	entities.PrintLog("SEND:\n", iso.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		er = iso.GoUnMarshal(st.Message)
		if er != nil {
			entities.PrintError(er.Error())
			return
		}
		// Override below:
		dataTrans.ResponseCode = iso.GetField(39)
		dataTrans.Msg = iso.GetField(48)
		entities.PrintLog("RECV:\n", iso.PrettyPrint())
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaction failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	// TODO: Begin Recycle apex transaction..
	if dataTrans.ResponseCode == constant.Success || (dataTrans.ResponseCode == constant.TransactedResponseGwLKM && dataTrans.Msg == helper.AlreadyReversed) {
		apexRepo, _ := apextransrepo.NewApexTransRepo()
		// Get Apx tx..
		var data entities.TransApx
		data, er = apexRepo.GetPrimaryTrxBelongToRecreateApx("TINTCR"+iso.GetField(11), "100", dataTrans.BankCode)
		if er != nil {
			entities.PrintError(fmt.Sprint(err.NoRecord, ", [TINTDB only]"))
		}
		// Create Apx tx..
		newTrx := data
		newTrx.Kode_trans = "290"
		newTrx.My_kode_trans = "200"
		newTrx.Keterangan = "Reversal " + data.Keterangan
		er = apexRepo.DuplicateTrxBelongToRecreateApx(newTrx)
		if er != nil {
			return er
		}
	}
	return nil
}

func (r *retransactionRepoImpl) ResendReversalGwlkmTransaction(dataTrans *entities.IsoMessageBody) (er error) {

	//TODO: Send ISO data to IP & Port GWLKM
	// ISO OBJ INIT
	iso, er := iso8583uParser.NewISO8583U()
	if er != nil {
		entities.PrintError("load package error", er.Error())
		return
	}

	// UNMARSHAL | RE-COMPOSE ISO
	er = iso.GoUnMarshal(dataTrans.Msg)
	if er != nil {
		entities.PrintError(er.Error())
		return
	}
	iso.SetMti(dataTrans.MTI)
	iso.SetField(3, "400700")

	// MARSHAL PROCS
	isoMsg, er := iso.GoMarshal()
	if er != nil {
		entities.PrintError(er.Error())
		return
	}

	// CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, er := repo.GetServeAddr(dataTrans.BankCode)
	if er != nil {
		entities.PrintError(er.Error())
	}

	// INIT TCP OBJ
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	entities.PrintLog("SEND:\n", iso.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		er = iso.GoUnMarshal(st.Message)
		if er != nil {
			entities.PrintError(er.Error())
			return
		}
		// Override below:
		dataTrans.ResponseCode = iso.GetField(39)
		dataTrans.Msg = iso.GetField(48)
		entities.PrintLog("RECV:\n", iso.PrettyPrint())
	} else {
		dataTrans.Msg = fmt.Sprint("reversal failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	return nil
}

func (r *retransactionRepoImpl) ResendReversalBeforeRecycleGwlkmTransaction(stan string) (er error) {

	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	var data entities.IsoMessageBody
	if data, er = dataRepo.GetData(stan); er != nil {
		return er
	}

	//TODO: Send ISO data to IP & Port GWLKM
	// ISO OBJ INIT
	iso, er := iso8583uParser.NewISO8583U()
	if er != nil {
		entities.PrintError("load package error", er.Error())
		return
	}

	// UNMARSHAL | RE-COMPOSE ISO
	er = iso.GoUnMarshal(data.Msg)
	if er != nil {
		entities.PrintError(er.Error())
		return
	}
	iso.SetMti(data.MTI)
	iso.SetField(3, "400700")

	// MARSHAL PROCS
	isoMsg, er := iso.GoMarshal()
	if er != nil {
		entities.PrintError(er.Error())
		return
	}

	// CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, er := repo.GetServeAddr(data.BankCode)
	if er != nil {
		entities.PrintError(er.Error())
	}

	// INIT TCP OBJ
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	entities.PrintLog("SEND:\n", iso.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		er = iso.GoUnMarshal(st.Message)
		if er != nil {
			entities.PrintError(er.Error())
			return
		}
		// Override below:
		data.ResponseCode = iso.GetField(39)
		data.Msg = iso.GetField(48)
		entities.PrintLog("RECV:\n", iso.PrettyPrint())

	} else {
		data.Msg = fmt.Sprint("Reversal failed: ", st.Message)
		return errors.New(data.Msg)
	}

	entities.PrintLog(data.Msg)
	return er
}

func (r *retransactionRepoImpl) RecycleSuspectRevBiller(dataTrans *entities.IsoMessageBody) (er error) {

	//TODO: Send ISO data to IP & Port GWLKM
	// ISO OBJ INIT
	iso, er := iso8583uParser.NewISO8583U()
	if er != nil {
		entities.PrintError("load package error", er.Error())
		return
	}

	// UNMARSHAL | RE-COMPOSE ISO
	er = iso.GoUnMarshal(dataTrans.Msg)
	if er != nil {
		entities.PrintError(er.Error())
		return
	}
	iso.SetMti(dataTrans.MTI)
	iso.SetField(3, "400700")
	iso.SetField(12, helper.GetCurrentDate())
	iso.SetField(104, dataTrans.ProductCode)

	// MARSHAL PROCS
	isoMsg, er := iso.GoMarshal()
	if er != nil {
		entities.PrintError(er.Error())
		return
	}

	// CORE ADDRS
	repo, _ := echanneltransrepo.NewEchannelTransRepo()
	coreAddr, er := repo.GetServeAddr(dataTrans.BankCode)
	if er != nil {
		entities.PrintError(er.Error())
	}

	// INIT TCP OBJ
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	entities.PrintLog("SEND:\n", iso.PrettyPrint())
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		er = iso.GoUnMarshal(st.Message)
		if er != nil {
			entities.PrintError(er.Error())
			return
		}
		// Override below:
		dataTrans.ResponseCode = iso.GetField(39)
		dataTrans.Msg = iso.GetField(48)
		entities.PrintLog("RECV:\n", iso.PrettyPrint())
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaction failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	// TODO: Begin Recycle apex transaction..
	if dataTrans.ResponseCode == constant.Success || (dataTrans.ResponseCode == constant.TransactedResponseGwLKM && dataTrans.Msg == helper.AlreadyReversed) {
		apexRepo, _ := apextransrepo.NewApexTransRepo()

		// TEXTCR Reversal..
		TEXTCR, er := apexRepo.GetPrimaryTrxBelongToRecreateApx("TEXTCR"+iso.GetField(11), "100", dataTrans.BankCode)
		if er != nil {
			entities.PrintError(fmt.Sprint(err.NoRecord, ", [TEXTDB only]"))
		}
		newTEXTCR := TEXTCR
		newTEXTCR.No_rekening = dataTrans.LKMSource
		newTEXTCR.Kode_trans = "290"
		newTEXTCR.My_kode_trans = "200"
		newTEXTCR.Keterangan = "Reversal " + TEXTCR.Keterangan
		if er = apexRepo.DuplicateTrxBelongToRecreateApx(newTEXTCR); er != nil {
			return er // berhenti disini jika sudah ada, kebawah tidak di eksekusi
		}

		// TEXTDB Reversal..
		TEXTDB, er := apexRepo.GetPrimaryTrxBelongToRecreateApx("TEXTDB"+iso.GetField(11), "200", dataTrans.BankCode)
		if er != nil {
			entities.PrintError(fmt.Sprint(err.NoRecord, ", [TEXTCR only]"))
		}
		newTEXTDB := TEXTDB
		newTEXTDB.No_rekening = dataTrans.LKMSource
		newTEXTDB.Kode_trans = "190"
		newTEXTDB.My_kode_trans = "100"
		newTEXTDB.Keterangan = "Reversal " + TEXTDB.Keterangan
		if er = apexRepo.DuplicateTrxBelongToRecreateApx(newTEXTDB); er != nil {
			return er
		}

		// TODO: Change EMS response code..
		trxSource, er := repo.GetOriginData(iso.GetField(11))
		if er != nil {
			return er
		}

		if er = repo.ChangeResponseCode("1100", trxSource.Stan, trxSource.Trans_id); er != nil {
			return er
		}

	}
	return nil
}
