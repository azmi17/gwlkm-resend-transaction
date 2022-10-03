package datatransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
)

func newDatatransRepoMysqlImpl(conn *sql.DB) DatatransRepo {
	return &DatatransRepoMysqlImpl{conn: conn}
}

type DatatransRepoMysqlImpl struct {
	conn *sql.DB
}

func (d *DatatransRepoMysqlImpl) GetData(stan string) (data entities.TransHisotry, er error) {
	row := d.conn.QueryRow(`SELECT 
		mti,
		bank_code,
		msg
		FROM trans_history WHERE stan= ?`, stan)

	er = row.Scan(
		&data.MTI,
		&data.BankCode,
		&data.Msg,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return data, nil
		} else {
			return data, errors.New(fmt.Sprint("error while get trans history : ", er.Error()))
		}
	}

	return
}
func (d *DatatransRepoMysqlImpl) GetServeAddr(bankCode string) (data entities.CoreAddr, er error) {
	row := d.conn.QueryRow(`SELECT 
		bank_code,
		ip_addr,
		ip_port
		FROM core_addrs WHERE bank_code= ?`, bankCode)

	er = row.Scan(
		&data.BankCode,
		&data.IPaddr,
		&data.TCPPort,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return data, nil
		} else {
			return data, errors.New(fmt.Sprint("error while get core addrs : ", er.Error()))
		}
	}

	return
}

// func (d *DatatransRepoMysqlImpl) GetData(stan string) (data entities.TransHisotry, er error) {
// 	row := d.conn.QueryRow(`SELECT
// 		mti,
// 		processing_code,
// 		amount,
// 		markup_total,
// 		profit_excluded,
// 		substr(tgl_trans_str,5) AS TmsTime,
// 		profit_included,
// 		stan,
// 		tgl_trans_str,
// 		substr(tgl_trans_str,5,4) AS dateTime,
// 		bank_code,
// 		msg
// 		FROM trans_history WHERE stan= ?`, stan)

// 	er = row.Scan(
// 		&data.MTI,
// 		&data.ProcessingCode,
// 		&data.Amount,
// 		&data.MarkupTotal,
// 		&data.ProfitEx,
// 		&data.TransmissionTime,
// 		&data.ProfitIn,
// 		&data.Stan,
// 		&data.DateTime,
// 		&data.Date,
// 		&data.BankCode,
// 		&data.ResponseMessage,
// 	)
// 	if er != nil {
// 		if er == sql.ErrNoRows {
// 			return data, nil
// 		} else {
// 			return data, errors.New(fmt.Sprint("error while get data : ", er.Error()))
// 		}
// 	}

// 	return
// }
