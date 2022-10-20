package apextransrepo

import (
	"database/sql"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"reflect"
	"testing"
	"time"

	"github.com/kpango/glg"
)

func GetConnectionApx() *sql.DB {
	dataSource := "root:azmic0ps@tcp(localhost:3317)/integrasi_apex_ems?parseTime=true"
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

func GetNextIdTabtrans() int {
	db := GetConnectionApx()
	userId := "1779"
	var transId int

	row := db.QueryRow("SELECT "+"ibs_get_next_id_with_userid(?) AS trans_id", userId)
	err := row.Scan(&transId)
	if err != nil {
		_ = glg.Log(err.Error())
	}

	return transId
}

func TestGetTabtransTx(t *testing.T) {
	db := GetConnectionApx()
	echanneltransrepo := newApexTransRepoMysqlImpl(db)
	data, err := echanneltransrepo.GetTxInfoApx("TKREDB830378018335")
	if err != nil {
		_ = glg.Log(err.Error())
	}

	v := reflect.ValueOf(data)
	typeOfData := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// fmt.Printf("Field: %s\t\t Value: %v\n", typeOfData.Field(i).Name, v.Field(i).Interface())
		fmt.Printf("%s => %v\n", typeOfData.Field(i).Name, v.Field(i).Interface())
	}

}

func TestCreateTabtransTx(t *testing.T) {
	db := GetConnectionApx()
	echanneltransrepo := newApexTransRepoMysqlImpl(db)

	// GET DATA
	var data entities.TransApx
	data, err := echanneltransrepo.GetTxInfoApx("TKREDB830378018335")
	if err != nil {
		_ = glg.Log(err.Error())
	}

	transId, err := echanneltransrepo.GetTransIdApx()
	if err != nil {
		_ = glg.Log(err.Error())
	}

	// CREATE DATA
	newTrx := data
	newTrx.Tabtrans_id = transId
	newTrx.Kuitansi = "S50RT4953021020"
	newTrx.Userid = 1779
	err = echanneltransrepo.DuplicatingTxApx(newTrx)
	if err != nil {
		_ = glg.Log(err.Error())
	}
	fmt.Println("Duplicating transaction succeeded..")

	// DELETE DATA
	err = echanneltransrepo.DeleteTxApx("TKREDB830378018335")
	if err != nil {
		_ = glg.Log(err.Error())
	}
	fmt.Println("Delete transaction succeeded..")

}

func TestInsertDummyTx(t *testing.T) {
	db := GetConnectionApx()
	echanneltransrepo := newApexTransRepoMysqlImpl(db)

	// GET DATA
	var data entities.TransApx
	data, err := echanneltransrepo.GetTxInfoApx("TKREDB109003590231")
	if err != nil {
		_ = glg.Log(err.Error())
	}

	// CREATE DATA 1
	transId, err := echanneltransrepo.GetTransIdApx()
	if err != nil {
		_ = glg.Log(err.Error())
	}
	newTrx1 := data
	newTrx1.Tabtrans_id = transId
	newTrx1.Kode_trans = "290"
	newTrx1.My_kode_trans = "200"
	newTrx1.Kuitansi = "S5100041274590"
	newTrx1.Pokok = 5000
	newTrx1.Userid = 1779
	newTrx1.Pay_biller_code = "400400"
	newTrx1.Pay_product_code = "S5"

	// CREATE DATA 2
	transId, err = echanneltransrepo.GetTransIdApx()
	if err != nil {
		_ = glg.Log(err.Error())
	}
	newTrx2 := data
	newTrx2.Tabtrans_id = transId
	newTrx2.Kode_trans = "190"
	newTrx2.My_kode_trans = "100"
	newTrx2.Kuitansi = "S5100041274590"
	newTrx2.Pokok = 5000
	newTrx2.Keterangan = "Reversal " + data.Keterangan
	newTrx2.Userid = 1779
	newTrx2.Pay_biller_code = "400400"
	newTrx2.Pay_product_code = "S5"

	err = echanneltransrepo.DuplicatingUnitTestTxApx(newTrx1, newTrx2)
	if err != nil {
		_ = glg.Log(err.Error())
	}

	fmt.Println("Duplicating transaction succeeded..")
}

func TestDeleteDummyTx(t *testing.T) {
	db := GetConnectionApx()
	echanneltransrepo := newApexTransRepoMysqlImpl(db)

	// DELETE DATA
	err := echanneltransrepo.DeleteTxApx("S5100041274590")
	if err != nil {
		_ = glg.Log(err.Error())
	}
	fmt.Println("Delete transaction succeeded..")

}
