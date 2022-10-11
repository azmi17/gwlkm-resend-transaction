package retransactionepo

import (
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/datatransrepo"

	iso8583uParser "github.com/randyardiansyah25/iso8583u/parser"
	"github.com/randyardiansyah25/libpkg/net/tcp"
)

func NewRetransactionRepo() RetransactionRepo {
	return &retransactionRepoImpl{}
}

type retransactionRepoImpl struct {
}

func (r *retransactionRepoImpl) RecycleTransaction(dataTrans *entities.MsgTransHistory) (err error) {

	//TODO: Send ISO data to ip & port gwlkm

	repo, _ := datatransrepo.NewDatatransRepo()

	// CALL CORE ADDR
	coreAddr, err := repo.GetServeAddr(dataTrans.BankCode)
	if err != nil {
		entities.PrintError(err.Error())
	}

	// TCP OBJ INIT..
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	st := client.Send(tcp.SetHeader(dataTrans.Msg, 4))
	//fmt.Println(st.Code, " : ", st.Message) // <= fmt.println change with ?

	// ISO OBJ INIT..
	isoUnMarshal, err := iso8583uParser.NewISO8583U()
	if err != nil {
		entities.PrintError("load package error", err.Error())
		return
	}

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			entities.PrintError(err.Error()) // <= err source ?
			return
		}
		// Override below:
		dataTrans.ResponseCode = isoUnMarshal.GetField(39)
		dataTrans.Msg = isoUnMarshal.GetField(48)
		isoUnMarshal.PrettyPrint()
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaciton failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	return nil
}

func (r *retransactionRepoImpl) RecycleReversedTransaction(dataTrans *entities.MsgTransHistory) (err error) {

	//TODO: Send ISO data to ip & port gwlkm

	// ISO OBJ INIT..
	isoUnMarshal, err := iso8583uParser.NewISO8583U()
	if err != nil {
		entities.PrintError("load package error", err.Error())
		return
	}

	//RE-COMPOSE ISO..
	err = isoUnMarshal.GoUnMarshal(dataTrans.Msg)
	if err != nil {
		entities.PrintError(err.Error()) // <= err source ?
		return
	}
	isoUnMarshal.SetMti(dataTrans.MTI)
	isoUnMarshal.SetField(3, isoUnMarshal.GetField(3))
	isoUnMarshal.SetField(4, isoUnMarshal.GetField(4))
	isoUnMarshal.SetField(5, isoUnMarshal.GetField(5))
	isoUnMarshal.SetField(6, isoUnMarshal.GetField(6))
	isoUnMarshal.SetField(7, isoUnMarshal.GetField(7))
	isoUnMarshal.SetField(8, isoUnMarshal.GetField(8))
	isoUnMarshal.SetField(11, dataTrans.Stan)
	isoUnMarshal.SetField(12, dataTrans.Date)
	isoUnMarshal.SetField(13, helper.GETMMDD[4:8])
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
	isoUnMarshal.SetField(104, isoUnMarshal.GetField(104))

	// MARSHAL PROCS..
	isoMsg, err := isoUnMarshal.GoMarshal()
	if err != nil {
		return
	}

	// CALL CORE ADDR
	repo, _ := datatransrepo.NewDatatransRepo()
	coreAddr, err := repo.GetServeAddr(dataTrans.BankCode)
	if err != nil {
		entities.PrintError(err.Error())
	}

	// TCP OBJ INIT..
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 45)
	st := client.Send(tcp.SetHeader(isoMsg, 4))

	// UNMARSHAL ISO FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			entities.PrintError(err.Error()) // <= err source ?
			return
		}
		// Override below:
		dataTrans.ResponseCode = isoUnMarshal.GetField(39)
		dataTrans.Msg = isoUnMarshal.GetField(48)
		isoUnMarshal.PrettyPrint()
	} else {
		dataTrans.Msg = fmt.Sprint("re-transaciton failed: ", st.Message)
		return errors.New(dataTrans.Msg)
	}

	return nil
}
