package echanneltransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/repository/apextransrepo"
	"gwlkm-resend-transaction/repository/constant"
	"time"
)

func newEchannelTransRepoMysqlImpl(echannelConn *sql.DB) EchannelTransRepo {
	return &EchannelTransRepoMysqlImpl{
		echannelDb: echannelConn,
	}
}

type EchannelTransRepoMysqlImpl struct {
	echannelDb *sql.DB
}

func (e *EchannelTransRepoMysqlImpl) GetData(stan string) (data entities.IsoMessageBody, er error) {
	row := e.echannelDb.QueryRow(`SELECT 
		mti,
		stan,
		tgl_trans_str,
		processing_code,
		bank_code,
		msg
		FROM trans_history WHERE stan= ? AND dc='d' LIMIT 1`, stan)

	er = row.Scan(
		&data.MTI,
		&data.Stan,
		&data.DateTime,
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

func (e *EchannelTransRepoMysqlImpl) GetServeAddr(bankCode string) (data entities.CoreAddrInfo, er error) {
	row := e.echannelDb.QueryRow(`SELECT 
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

func (e *EchannelTransRepoMysqlImpl) GetOriginData(stan string) (data entities.TransHistory, er error) {

	row := e.echannelDb.QueryRow(`SELECT 
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

func (e *EchannelTransRepoMysqlImpl) DuplicatingData(copy entities.TransHistory) (er error) {

	stmt, er := e.echannelDb.Prepare(`INSERT INTO trans_history(
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
		dataRepo, _ := NewEchannelTransRepo()
		addStanReference := entities.StanReferences{
			Trans_id: int(lastId),
			Ref_Stan: copy.Ref_Stan,
			Stan:     copy.Stan,
		}
		er = dataRepo.AddStanReference(addStanReference)
		if er != nil {
			return er
		}
		er = dataRepo.ChangeResponseCode(constant.Resend, copy.Ref_Stan, copy.Trans_id)
		if er != nil {
			return err.InternalServiceError
		}

		// TODO: Recycle apex transaction..
		dataRepoApex, _ := apextransrepo.NewApexTransRepo()

		// Get Apx tx..
		var data entities.TransApx
		data, er = dataRepoApex.GetTabtransTxInfoApx(copy.Product_Code+copy.Ref_Stan, copy.Bank_Code)
		if er != nil {
			return er
		}

		// Create Apx tx..
		newTrx := data
		newTrx.Tgl_trans = time.Now()
		newTrx.Kuitansi = copy.Product_Code + copy.Stan
		newTrx.Userid = helper.GetUserIDApp()
		er = dataRepoApex.DuplicatingTxApx(newTrx)
		if er != nil {
			return er
		}

		// Delete Apx tx..
		er = dataRepoApex.DeleteTxApx(copy.Product_Code+copy.Ref_Stan, copy.Bank_Code)
		if er != nil {
			return er
		}

		return nil
	}
}

func (e *EchannelTransRepoMysqlImpl) ChangeResponseCode(rc, stan string, transId int) (er error) {
	var stmt *sql.Stmt
	if transId == 0 {
		// Search STAN is exist or not (?)
		dataRepo, _ := NewEchannelTransRepo()
		_, er = dataRepo.GetRetransTxInfo(stan)
		if er != nil {
			return err.NoRecord
		}

		stmt, er = e.echannelDb.Prepare(`UPDATE trans_history SET response_code = ? WHERE stan = ? AND dc='d'`)
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
		stmt, er = e.echannelDb.Prepare(`UPDATE trans_history SET response_code = ? WHERE stan = ? AND trans_id = ? AND dc='d'`)
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

func (e *EchannelTransRepoMysqlImpl) AddStanReference(reference entities.StanReferences) (er error) {

	stmt, er := e.echannelDb.Prepare(`INSERT INTO stan_ref_retrans(trans_id, ref_stan, stan) VALUES(?,?,?)`)
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

func (e *EchannelTransRepoMysqlImpl) GetRetransTxInfo(stan string) (txInfo web.RetransTxInfo, er error) {
	row := e.echannelDb.QueryRow(`SELECT 
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
		&txInfo.Iso_Msg_Resp,
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

func (e *EchannelTransRepoMysqlImpl) UpdateIsoMsg(isoMsg, stan string) (er error) {

	dataRepo, _ := NewEchannelTransRepo()

	_, er = dataRepo.GetRetransTxInfo(stan)
	if er != nil {
		return err.NoRecord
	}

	stmt, er := e.echannelDb.Prepare(`UPDATE trans_history SET msg_response = ? WHERE stan = ? AND dc='d'`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare update message response: ", er.Error()))
	}

	defer func() {
		_ = stmt.Close()
	}()

	if _, er = stmt.Exec(isoMsg, stan); er != nil {
		return errors.New(fmt.Sprint("error while update message response: ", er.Error()))
	}

	return nil
}
