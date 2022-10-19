package datatransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/constant"
)

func newDatatransRepoMysqlImpl(conn1, conn2 *sql.DB) DatatransRepo {
	return &DatatransRepoMysqlImpl{
		db1: conn1,
		db2: conn2,
	}
}

type DatatransRepoMysqlImpl struct {
	db1, db2 *sql.DB
}

func (d *DatatransRepoMysqlImpl) GetData(stan string) (data entities.MsgTransHistory, er error) {
	row := d.db1.QueryRow(`SELECT 
		mti,
		processing_code,
		bank_code,
		msg
		FROM trans_history WHERE stan= ? AND dc='d' LIMIT 1`, stan)

	er = row.Scan(
		&data.MTI,
		&data.ProcessingCode,
		&data.BankCode,
		&data.Msg,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return data, err.NoRecord
		} else {
			return data, errors.New(fmt.Sprint("error while get data: ", er.Error()))
		}
	}

	return
}

func (d *DatatransRepoMysqlImpl) GetServeAddr(bankCode string) (data entities.CoreAddr, er error) {
	row := d.db1.QueryRow(`SELECT 
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

	row := d.db1.QueryRow(`SELECT 
		th.trans_id,
		th.stan,
		COALESCE(sr.ref_stan, '') as ref_stan,
		th.tgl_trans_str,
		th.bank_code,
		th.rek_id,
		th.mti,
		th.processing_code,
		th.biller_code,
		th.product_code,
		th.subscriber_id,
		th.dc,
		th.response_code,
		th.amount,
		th.qty,
		th.profit_included,
		th.profit_excluded,
		th.profit_share_biller,
		th.profit_share_aggregator,
		th.profit_share_bank,
		th.markup_total,
		th.markup_share_aggregator,
		th.markup_share_bank,
		th.msg,
		th.msg_response,
		th.bit39_bit48_hulu,
		th.saldo_before_trans,
		th.keterangan,
		th.ref,
		th.synced_ibs_core,
		th.synced_ibs_core_description,
		th.bris_original_data,
		th.gateway_id,
		th.id_user,
		th.id_raw,
		th.advice_count,
		th.status_id,
		th.nohp_notif,
		th.score,
		th.no_hp_alternatif,
		th.inc_notif_status,
		th.fee_rek_induk
	FROM trans_history AS th 
	LEFT JOIN stan_ref_retrans AS sr ON(th.trans_id=sr.trans_id) WHERE th.stan= ? AND dc='d' LIMIT 1`, stan)

	er = row.Scan(
		&data.Trans_id,
		&data.Stan,
		&data.Ref_Stan,
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
			return data, err.NoRecord
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

func (d *DatatransRepoMysqlImpl) DuplicatingData(copy entities.TransHistory) (er error) {

	stmt, er := d.db1.Prepare(`INSERT INTO trans_history(
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
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare duplicating transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if stmt, er := stmt.Exec(
		copy.Stan,
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
		return errors.New(fmt.Sprint("error while duplicating transaction: ", er.Error()))
	} else {

		// Get Last insert ID
		lastId, txErr := stmt.LastInsertId()
		if txErr != nil {
			return errors.New(fmt.Sprint("error while get last insert id - add stan reference: ", txErr.Error()))
		}

		// below is to add ref_stan to stan_ref_retrans and change response_code belonging to the first record..
		dataRepo, _ := NewDatatransRepo()
		addStanReference := entities.StanReference{
			Trans_id: int(lastId),
			Ref_Stan: copy.Ref_Stan,
			Stan:     copy.Stan,
		}
		er = dataRepo.AddStanReference(addStanReference)
		if er != nil {
			return er
		}
		er = dataRepo.ChangeResponseCode(constant.Failed, copy.Ref_Stan, copy.Trans_id)
		if er != nil {
			return err.InternalServiceError
		}

		//TODO: Recycle apex transaction..
		// Get Trans ID
		transId, err := dataRepo.GetTransIdApx()
		if er != nil {
			return err
		}

		// Get Apx tx..
		var data entities.TransApx
		data, er = dataRepo.GetTxInfoApx(copy.Product_Code + copy.Ref_Stan)
		if er != nil {
			return er
		}

		// Create Apx tx..
		newTrx := data
		newTrx.Tabtrans_id = transId
		newTrx.Kuitansi = copy.Product_Code + copy.Stan
		newTrx.Userid = helper.GetUserIDApp()
		err = dataRepo.DuplicatingTxApx(newTrx)
		if err != nil {
			return er
		}

		// Delete Apx tx..
		err = dataRepo.DeleteTxApx(copy.Product_Code + copy.Ref_Stan)
		if err != nil {
			return er
		}

		return nil
	}
}

func (d *DatatransRepoMysqlImpl) ChangeResponseCode(rc, stan string, transId int) (er error) {
	var stmt *sql.Stmt
	if transId == 0 {

		// Search STAN is exist or not (?)
		dataRepo, _ := NewDatatransRepo()
		_, er = dataRepo.GetRetransTxInfo(stan)
		if er != nil {
			return err.NoRecord
		}

		// Upd procs..
		stmt, er = d.db1.Prepare(`UPDATE trans_history SET response_code = ? WHERE stan = ? AND dc='d'`)
		if er != nil {
			return errors.New(fmt.Sprint("error while prepare update response code: ", er.Error()))
		}

		defer func() {
			_ = stmt.Close()
		}()

		if _, er = stmt.Exec(rc, stan); er != nil {
			return errors.New(fmt.Sprint("error while update response code: ", er.Error()))
		}
	} else {
		stmt, er = d.db1.Prepare(`UPDATE trans_history SET response_code = ? WHERE stan = ? AND trans_id = ? AND dc='d'`)
		if er != nil {
			return errors.New(fmt.Sprint("error while prepare update retrans response code: ", er.Error()))
		}

		defer func() {
			_ = stmt.Close()
		}()

		if _, er = stmt.Exec(rc, stan, transId); er != nil {
			return errors.New(fmt.Sprint("error while update response code: ", er.Error()))
		}
	}
	return nil
}

func (d *DatatransRepoMysqlImpl) AddStanReference(reference entities.StanReference) (er error) {

	stmt, er := d.db1.Prepare(`INSERT INTO stan_ref_retrans(trans_id, ref_stan, stan) VALUES(?,?,?)`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare add stan reference: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er = stmt.Exec(
		reference.Trans_id,
		reference.Ref_Stan,
		reference.Stan,
	); er != nil {
		return errors.New(fmt.Sprint("error while add stan reference: ", er.Error()))
	} else {
		return nil
	}
}

func (d *DatatransRepoMysqlImpl) GetRetransTxInfo(stan string) (txInfo entities.RetransTxInfo, er error) {
	row := d.db1.QueryRow(`SELECT 
		th.trans_id,
		th.subscriber_id,
		th.stan,
		COALESCE(sr.ref_stan, '') as ref_stan,
		DATE_FORMAT(th.tgl_trans_str, "%d/%m/%Y") AS tgl_trans,
		th.response_code,
		th.bank_code,
		th.rek_id,
		th.amount,
		th.ref AS kuitansi,
		th.msg_response
	FROM trans_history AS th 
	LEFT JOIN stan_ref_retrans AS sr ON(th.trans_id=sr.trans_id) WHERE th.stan= ? AND dc='d' LIMIT 1`, stan)
	er = row.Scan(
		&txInfo.Trans_Id,
		&txInfo.Idpel,
		&txInfo.Stan,
		&txInfo.Ref_Stan,
		&txInfo.Tgl_Trans_Str,
		&txInfo.Response_Code,
		&txInfo.BankCode,
		&txInfo.RekID,
		&txInfo.Amount,
		&txInfo.Kuitansi,
		&txInfo.Iso_Msg,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return txInfo, err.NoRecord
		} else {
			return txInfo, errors.New(fmt.Sprint("error while get retrans data: ", er.Error()))
		}
	}
	return
}

func (d *DatatransRepoMysqlImpl) GetTransIdApx() (transId int, er error) {

	userId := helper.GetUserIDApp()
	row := d.db2.QueryRow(`SELECT ibs_get_next_id_with_userid(?) AS trans_id`, userId)
	er = row.Scan(
		&transId,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return transId, err.NoRecord
		} else {
			return transId, errors.New(fmt.Sprint("error while get data: ", er.Error()))
		}
	}
	return transId, nil
}

func (d *DatatransRepoMysqlImpl) GetTxInfoApx(kuitansi string) (transApx entities.TransApx, er error) {
	row := d.db2.QueryRow(`SELECT 
		tabtrans_id,
		tgl_trans,
		no_rekening,
		kode_trans,
		my_kode_trans,
		pokok,
		kuitansi,
		userid,
		keterangan,
		verifikasi,
		tob,
		sandi_trans,
		posted_to_gl,
		kode_kantor,
		jam,
		tgl_real_trans,
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	FROM tabtrans WHERE kuitansi= ? AND my_kode_trans='200' LIMIT 1`, kuitansi)
	er = row.Scan(
		&transApx.Tabtrans_id,
		&transApx.Tgl_trans,
		&transApx.No_rekening,
		&transApx.Kode_trans,
		&transApx.My_kode_trans,
		&transApx.Pokok,
		&transApx.Kuitansi,
		&transApx.Userid,
		&transApx.Keterangan,
		&transApx.Verifikasi,
		&transApx.Tob,
		&transApx.Sandi_trans,
		&transApx.Posted_to_gl,
		&transApx.Kode_kantor,
		&transApx.Jam,
		&transApx.Tgl_real_trans,
		&transApx.Pay_lkm_source,
		&transApx.Pay_lkm_norek,
		&transApx.Pay_idpel,
		&transApx.Pay_biller_code,
		&transApx.Pay_product_code,
	)
	if er != nil {
		if er == sql.ErrNoRows {
			return transApx, err.NoRecord
		} else {
			return transApx, errors.New(fmt.Sprint("error while get tabtrans data: ", er.Error()))
		}
	}
	return
}

func (d *DatatransRepoMysqlImpl) DuplicatingTxApx(copy entities.TransApx) (er error) {

	stmt, er := d.db2.Prepare(`INSERT INTO tabtrans(
		tabtrans_id,
		tgl_trans,
		no_rekening,
		kode_trans,
		my_kode_trans,
		pokok,
		kuitansi,
		userid,
		keterangan,
		verifikasi,
		tob,
		sandi_trans,
		posted_to_gl,
		kode_kantor,
		jam,
		tgl_real_trans,
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare add tabtrans transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er := stmt.Exec(
		copy.Tabtrans_id,
		copy.Tgl_trans,
		copy.No_rekening,
		copy.Kode_trans,
		copy.My_kode_trans,
		copy.Pokok,
		copy.Kuitansi,
		copy.Userid,
		copy.Keterangan,
		copy.Verifikasi,
		copy.Tob,
		copy.Sandi_trans,
		copy.Posted_to_gl,
		copy.Kode_kantor,
		copy.Jam,
		copy.Tgl_real_trans,
		copy.Pay_lkm_source,
		copy.Pay_lkm_norek,
		copy.Pay_idpel,
		copy.Pay_biller_code,
		copy.Pay_product_code); er != nil {
		return errors.New(fmt.Sprint("error while add tabtrans transaction: ", er.Error()))
	} else {
		return nil
	}
}

func (d *DatatransRepoMysqlImpl) DeleteTxApx(kuitansi string) (er error) {

	stmt, er := d.db2.Prepare("DELETE FROM tabtrans WHERE kuitansi = ?")
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare delete tabtrans transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er := stmt.Exec(kuitansi); er != nil {
		return errors.New(fmt.Sprint("error while delete tabtrans transaction: ", er.Error()))
	}

	return nil
}

func (d *DatatransRepoMysqlImpl) RecycleTxApx(kuitansi string) (er error) {
	dataRepo, _ := NewDatatransRepo()

	// Get Trans ID
	transId, err := dataRepo.GetTransIdApx()
	if er != nil {
		return err
	}

	// Get DATA
	var data entities.TransApx
	data, er = dataRepo.GetTxInfoApx(kuitansi)
	if er != nil {
		return er
	}

	// CREATE DATA
	newTrx := data
	newTrx.Tabtrans_id = transId
	newTrx.Kuitansi = "S50RT4953021020"
	newTrx.Userid = 1779
	err = dataRepo.DuplicatingTxApx(newTrx)
	if err != nil {
		return er
	}

	// DELETE DATA
	err = dataRepo.DeleteTxApx(kuitansi)
	if err != nil {
		return er
	}

	return nil
}
