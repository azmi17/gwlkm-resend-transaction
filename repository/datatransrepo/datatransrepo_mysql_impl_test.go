package datatransrepo

import (
	"database/sql"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/helper"
	"log"
	"testing"
	"time"

	"github.com/randyardiansyah25/libpkg/net/tcp"

	iso8583uParser "github.com/randyardiansyah25/iso8583u/parser"
)

func GetConnection() *sql.DB {
	dataSource := "root:@tcp(localhost:3317)/echannel?parseTime=true"
	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}

func TestGetData(t *testing.T) {

	// CALL BANK DATA REPO
	db := GetConnection()
	dataTransRepo := newDatatransRepoMysqlImpl(db)
	data, err := dataTransRepo.GetData("100041274590")
	if err != nil {
		log.Fatal(err.Error())
	}
	coreAddr, _ := dataTransRepo.GetServeAddr(data.BankCode)

	// CALL RE-TRANSACTION REPO
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 30)
	st := client.Send(tcp.SetHeader(data.Msg, 4))
	fmt.Println(st.Code, " : ", st.Message)

	// ISO OBJ INIT..
	isoUnMarshal, _ := iso8583uParser.NewISO8583U()
	if err != nil {
		fmt.Println("load package error", err.Error())
		return
	}

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log("Message :\n", st.Message)
		t.Log("Parse : \n", isoUnMarshal.PrettyPrint())

		// rc := isoUnMarshal.GetField(39)
		// if rc != "0000" {

		// }
	}

}

func TestMarshalISO(t *testing.T) {

	// CALL BANK DATA REPO
	db := GetConnection()
	dataTransRepo := newDatatransRepoMysqlImpl(db)
	data, err := dataTransRepo.GetData("100041274590")
	if err != nil {
		log.Fatal(err.Error())
	}

	isoUnMarshal, _ := iso8583uParser.NewISO8583U()
	if err != nil {
		fmt.Println("load package error", err.Error())
		return
	}
	err = isoUnMarshal.GoUnMarshal(data.Msg)
	if err != nil {
		t.Error(err.Error())
		return
	}

	t.Log("Message :\n", data.Msg)
	t.Log("Parse : \n", isoUnMarshal.PrettyPrint())
	t.Log("Get Amount : ", isoUnMarshal.GetField(4))

}

func TestReversedData(t *testing.T) {

	// CALL BANK DATA REPO
	db := GetConnection()
	dataTransRepo := newDatatransRepoMysqlImpl(db)

	// CALL REVERSED DATA
	reversedData, err := dataTransRepo.GetReversedData("100041274590")
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println(reversedData.Time_Stamp)
	// for _, num := range data {

	// }

	// RE-COMPOSE
	newTrx := entities.TransHistory{}
	newTrx.Stan = "RT" + reversedData.Stan[2:12]   // thinking abt it
	newTrx.Tgl_Trans_Str = helper.GetCurrentDate() // thinking abt it
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
	newTrx.Time_Stamp = time.Now()
	newTrx.Saldo_Before_Trans = reversedData.Saldo_Before_Trans
	newTrx.Keterangan = reversedData.Keterangan
	newTrx.Ref = reversedData.Ref
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
	trx, err := dataTransRepo.DuplicatingData(newTrx)
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println("insert new transaction succeeded")

	// CALL CHANGE RESPONSE CODE
	err = dataTransRepo.ChangeRcOnReversedData(reversedData.Stan)
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println("Update reversed transaction succeeded")

	// CALL CORE-ADDRS
	coreAddr, _ := dataTransRepo.GetServeAddr(trx.Bank_Code)

	// CALL TCP AND SEND ISO MSG
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 30)
	st := client.Send(tcp.SetHeader(trx.Msg, 4))
	fmt.Println(st.Code, " : ", st.Message)

	// ISO OBJ INIT..
	isoUnMarshal, _ := iso8583uParser.NewISO8583U()
	if err != nil {
		fmt.Println("load package error", err.Error())
		return
	}

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log("Message :\n", st.Message)
		t.Log("Parse : \n", isoUnMarshal.PrettyPrint())
	}
}

func TestReversedDataTwo(t *testing.T) {

	// CALL BANK DATA REPO
	db := GetConnection()
	dataTransRepo := newDatatransRepoMysqlImpl(db)

	// CALL REVERSED DATA
	reversedData, err := dataTransRepo.GetReversedData("100041274590")
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println(reversedData.Time_Stamp)
	// for _, num := range data {

	// }

	// RE-COMPOSE
	newTrx := entities.TransHistory{}
	newTrx.Stan = "RT" + reversedData.Stan[2:12]   // thinking abt it
	newTrx.Tgl_Trans_Str = helper.GetCurrentDate() // thinking abt it
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
	newTrx.Time_Stamp = time.Now()
	newTrx.Saldo_Before_Trans = reversedData.Saldo_Before_Trans
	newTrx.Keterangan = reversedData.Keterangan
	newTrx.Ref = reversedData.Ref
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
	trx, err := dataTransRepo.DuplicatingData(newTrx)
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println("insert new transaction succeeded")

	// CALL CHANGE RESPONSE CODE
	err = dataTransRepo.ChangeRcOnReversedData(reversedData.Stan)
	if err != nil {
		log.Fatal(err.Error())
	}
	// fmt.Println("Update reversed transaction succeeded")

	isoUnMarshal, _ := iso8583uParser.NewISO8583U()
	if err != nil {
		fmt.Println("load package error", err.Error())
		return
	}
	err = isoUnMarshal.GoUnMarshal(trx.Msg)
	if err != nil {
		t.Error(err.Error())
		return
	}

	// t.Log("Message :\n", trx.Msg)
	// t.Log("Parse : \n", isoUnMarshal.PrettyPrint())

	//Marshal Procs..
	isoUnMarshal.SetMti(trx.Mti)
	isoUnMarshal.SetField(3, isoUnMarshal.GetField(3))
	isoUnMarshal.SetField(4, isoUnMarshal.GetField(4))
	isoUnMarshal.SetField(5, isoUnMarshal.GetField(5))
	isoUnMarshal.SetField(6, isoUnMarshal.GetField(6))
	isoUnMarshal.SetField(7, isoUnMarshal.GetField(7))
	isoUnMarshal.SetField(8, isoUnMarshal.GetField(8))
	isoUnMarshal.SetField(11, trx.Stan)
	isoUnMarshal.SetField(12, trx.Tgl_Trans_Str)
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
	isoUnMarshal.SetField(104, isoUnMarshal.GetField(104))
	isoMsg, err := isoUnMarshal.GoMarshal()
	if err != nil {
		t.Error(err.Error())
		return
	}
	// t.Log("Set Fields : \n", isoUnMarshal.PrettyPrint())
	// t.Log("Result : ", isoMsg)

	// CALL CORE-ADDRS
	coreAddr, _ := dataTransRepo.GetServeAddr(trx.Bank_Code)

	// CALL TCP AND SEND ISO MSG
	client := tcp.NewTCPClient(coreAddr.IPaddr, coreAddr.TCPPort, 30)
	st := client.Send(tcp.SetHeader(isoMsg, 4))
	fmt.Println(st.Code, " : ", st.Message)

	// UNMARSHAL PROCS FROM SENDER
	if st.Code == tcp.CONNOK {
		err = isoUnMarshal.GoUnMarshal(st.Message)
		if err != nil {
			t.Error(err.Error())
			return
		}
		t.Log("Message :\n", st.Message)
		t.Log("Parse : \n", isoUnMarshal.PrettyPrint())
	}
}

func TestGetCurrentDate(t *testing.T) {
	now := helper.GetCurrentDate()
	fmt.Println(now)
}
