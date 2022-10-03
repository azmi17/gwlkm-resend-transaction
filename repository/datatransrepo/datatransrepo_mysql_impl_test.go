package datatransrepo

import (
	"database/sql"
	"errors"
	"fmt"
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
	} else {
		errors.New(fmt.Sprint("re-transaciton failed: ", st.Message))
	}

}
