package datatransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/repository/constant"
)

func newDatatransRepoMysqlImpl(conn *sql.DB) DatatransRepo {
	return &DatatransRepoMysqlImpl{conn: conn}
}

type DatatransRepoMysqlImpl struct {
	conn *sql.DB
}

func (d *DatatransRepoMysqlImpl) GetData(stan string) (data entities.MsgTransHistory, er error) {
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
			return data, errors.New(fmt.Sprint("error while get data: ", er.Error()))
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
			return data, errors.New(fmt.Sprint("error while get core addrs: ", er.Error()))
		}
	}

	return
}

func (d *DatatransRepoMysqlImpl) GetReversedData(stan string) (data entities.TransHistory, er error) {
	row := d.conn.QueryRow(`SELECT 
		stan,
		tgl_trans_str,
		bank_code,
		rek_id,
		mti,
		processing_code,
		biller_code,
		product_code,
		subscriber_id,
		dc,
		response_code,
		amount,
		qty,
		profit_included,
		profit_excluded,
		profit_share_biller,
		profit_share_aggregator,
		profit_share_bank,
		markup_total,
		markup_share_aggregator,
		markup_share_bank,
		msg,
		msg_response,
		bit39_bit48_hulu,
		saldo_before_trans,
		keterangan,
		ref,
		synced_ibs_core,
		synced_ibs_core_description,
		bris_original_data,
		gateway_id,
		id_user,
		id_raw,
		advice_count,
		status_id,
		nohp_notif,
		score,
		no_hp_alternatif,
		inc_notif_status,
		fee_rek_induk
	FROM trans_history WHERE stan= ?`, stan)

	er = row.Scan(
		&data.Stan,
		&data.Tgl_Trans_Str,
		&data.Bank_Code,
		&data.Rek_Id,
		&data.Mti,
		&data.Processing_Code,
		&data.Biller_Code,
		&data.Product_Code,
		&data.Subscriber_Id,
		&data.Dc,
		&data.Response_Code,
		&data.Amount,
		&data.Qty,
		&data.Profit_Included,
		&data.Profit_Excluded,
		&data.Profit_Share_Biller,
		&data.Profit_Share_Aggregator,
		&data.Profit_Share_Bank,
		&data.Markup_Total,
		&data.Markup_Share_Aggregator,
		&data.Markup_Share_Bank,
		&data.Msg,
		&data.Msg_Response,
		&data.Bit39_Bit48_Hulu,
		&constant.SQLsaldo_before_trans,
		&data.Keterangan,
		&data.Ref,
		&data.Synced_Ibs_Core,
		&constant.SQLsynced_ibs_core_description,
		&data.Bris_Original_Data,
		&data.Gateway_Id,
		&constant.SQLid_user,
		&constant.SQLid_raw,
		&data.Advice_Count,
		&data.Status_Id,
		&data.Nohp_Notif,
		&data.Score,
		&constant.SQLno_hp_alternatif,
		&data.Inc_Notif_Status,
		&data.Fee_Rek_Induk,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return data, nil
		} else {
			return data, errors.New(fmt.Sprint("error while get reversed data: ", er.Error()))
		}
	}
	data.Saldo_Before_Trans = constant.SQLsaldo_before_trans.GetInt()
	data.Synced_Ibs_Core_Description = constant.SQLsynced_ibs_core_description.String
	data.Id_User = constant.SQLid_user.GetInt()
	data.Id_Raw = constant.SQLid_raw.GetInt()
	data.No_Hp_Alternatif = constant.SQLno_hp_alternatif.String

	return
}

func (d *DatatransRepoMysqlImpl) DuplicatingData(copy entities.TransHistory) (data entities.TransHistory, er error) {

	stmt, er := d.conn.Prepare(`INSERT INTO trans_history(
		stan,
		ref_stan,
		tgl_trans_str,
		bank_code,
		rek_id,
		mti,
		processing_code,
		biller_code,
		product_code,
		subscriber_id,
		dc,
		response_code,
		amount,
		qty,
		profit_included,
		profit_excluded,
		profit_share_biller,
		profit_share_aggregator,
		profit_share_bank,
		markup_total,
		markup_share_aggregator,
		markup_share_bank,
		msg,
		msg_response,
		bit39_bit48_hulu,
		saldo_before_trans,
		keterangan,
		ref,
		synced_ibs_core,
		synced_ibs_core_description,
		bris_original_data,
		gateway_id,
		id_user,
		id_raw,
		advice_count,
		status_id,
		nohp_notif,
		score,
		no_hp_alternatif,
		inc_notif_status,
		fee_rek_induk
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if er != nil {
		return data, errors.New(fmt.Sprint("error while prepare duplicating transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er = stmt.Exec(
		copy.Stan,
		copy.Ref_Stan,
		copy.Tgl_Trans_Str,
		copy.Bank_Code,
		copy.Rek_Id,
		copy.Mti,
		copy.Processing_Code,
		copy.Biller_Code,
		copy.Product_Code,
		copy.Subscriber_Id,
		copy.Dc,
		copy.Response_Code,
		copy.Amount,
		copy.Qty,
		copy.Profit_Included,
		copy.Profit_Excluded,
		copy.Profit_Share_Biller,
		copy.Profit_Share_Aggregator,
		copy.Profit_Share_Bank,
		copy.Markup_Total,
		copy.Markup_Share_Aggregator,
		copy.Markup_Share_Bank,
		copy.Msg,
		copy.Msg_Response,
		copy.Bit39_Bit48_Hulu,
		copy.Saldo_Before_Trans,
		copy.Keterangan,
		copy.Ref,
		copy.Synced_Ibs_Core,
		copy.Synced_Ibs_Core_Description,
		copy.Bris_Original_Data,
		copy.Gateway_Id,
		copy.Id_User,
		copy.Id_Raw,
		copy.Advice_Count,
		copy.Status_Id,
		copy.Nohp_Notif,
		copy.Score,
		copy.No_Hp_Alternatif,
		copy.Inc_Notif_Status,
		copy.Fee_Rek_Induk,
	); er != nil {
		return data, errors.New(fmt.Sprint("error while duplicating transaction: ", er.Error()))
	} else {
		return copy, nil
	}
}

func (d *DatatransRepoMysqlImpl) ChangeRcOnReversedData(rc, stan string) (er error) {
	stmt, er := d.conn.Prepare("UPDATE trans_history SET response_code = ? WHERE stan = ?")
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare update response code: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er = stmt.Exec(rc, stan); er != nil {
		return errors.New(fmt.Sprint("error while update response code: ", er.Error()))
	}

	return nil
}

func (d *DatatransRepoMysqlImpl) RollbackDuplicate(stan string) (er error) {
	stmt, er := d.conn.Prepare("DELETE FROM trans_history WHERE ref_stan = ?")
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare delete duplicated transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er := stmt.Exec(stan); er != nil {
		return errors.New(fmt.Sprint("error while delete duplicated transaction: ", er.Error()))
	}

	return nil
}
