package apextransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/helper"
	"sync"
)

func newApexTransRepoMysqlImpl(apexConn *sql.DB) ApexTransRepo {
	return &ApexTransRepoMysqlImpl{
		apexDb: apexConn,
	}
}

type ApexTransRepoMysqlImpl struct {
	apexDb *sql.DB
}

func (a *ApexTransRepoMysqlImpl) GetTransIdApx() (transId int, er error) {

	userId := helper.GetUserIDApp()
	row := a.apexDb.QueryRow(`SELECT ibs_get_next_id_with_userid(?) AS trans_id`, userId)
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

func (a *ApexTransRepoMysqlImpl) GetTabtransTxInfoApx(kuitansi, bankCode string) (transApx entities.TransApx, er error) {
	row := a.apexDb.QueryRow(`SELECT 
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
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	FROM tabtrans WHERE kuitansi= ? AND pay_lkm_source = ? AND my_kode_trans='200' LIMIT 1`, kuitansi, bankCode)
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

func (a *ApexTransRepoMysqlImpl) DuplicatingTxApx(copy ...entities.TransApx) (er error) {

	apexTransRepo, _ := NewApexTransRepo()

	stmt, er := a.apexDb.Prepare(`INSERT INTO tabtrans(
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
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare add tabtrans transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	for _, item := range copy {

		// Get Trans ID
		tabtransId, err := apexTransRepo.GetTransIdApx()
		if er != nil {
			return err
		}

		if _, er := stmt.Exec(
			tabtransId,
			item.Tgl_trans,
			item.No_rekening,
			item.Kode_trans,
			item.My_kode_trans,
			item.Pokok,
			item.Kuitansi,
			item.Userid,
			item.Keterangan,
			item.Verifikasi,
			item.Tob,
			item.Sandi_trans,
			item.Posted_to_gl,
			item.Kode_kantor,
			item.Jam,
			item.Pay_lkm_source,
			item.Pay_lkm_norek,
			item.Pay_idpel,
			item.Pay_biller_code,
			item.Pay_product_code); er != nil {
			return errors.New(fmt.Sprint("error while add tabtrans transaction: ", er.Error()))
		}
	}
	return nil
}

func (a *ApexTransRepoMysqlImpl) DuplicatingUnitTestTxApx(copy ...entities.TransApx) (er error) {

	stmt, er := a.apexDb.Prepare(`INSERT INTO tabtrans(
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
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare add tabtrans transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	for _, item := range copy {
		if _, er := stmt.Exec(
			item.Tabtrans_id,
			item.Tgl_trans,
			item.No_rekening,
			item.Kode_trans,
			item.My_kode_trans,
			item.Pokok,
			item.Kuitansi,
			item.Userid,
			item.Keterangan,
			item.Verifikasi,
			item.Tob,
			item.Sandi_trans,
			item.Posted_to_gl,
			item.Kode_kantor,
			item.Jam,
			item.Pay_lkm_source,
			item.Pay_lkm_norek,
			item.Pay_idpel,
			item.Pay_biller_code,
			item.Pay_product_code); er != nil {
			return errors.New(fmt.Sprint("error while add tabtrans transaction: ", er.Error()))
		}
	}
	return nil
}

func (a *ApexTransRepoMysqlImpl) DeleteTxApx(kuitansi, bankCode string) (er error) {

	stmt, er := a.apexDb.Prepare("DELETE FROM tabtrans WHERE kuitansi = ? AND pay_lkm_source = ?")
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare delete tabtrans transaction: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()

	if _, er := stmt.Exec(kuitansi, bankCode); er != nil {
		return errors.New(fmt.Sprint("error while delete tabtrans transaction: ", er.Error()))
	}

	return nil
}

func (a *ApexTransRepoMysqlImpl) GetTabtransListApx(kuitansi string) (listTx []web.TabtransInfoApx, er error) {
	rows, er := a.apexDb.Query(`SELECT
		DATE_FORMAT(tgl_trans, "%d/%m/%Y") AS tgl_trans,
		no_rekening,
		kode_trans,
		pay_idpel,
		kuitansi,
		pokok,
		keterangan,
		pay_product_code,
		pay_biller_code,
		userid
	FROM tabtrans WHERE kuitansi = ?`, kuitansi)
	if er != nil {
		return listTx, er
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var tabtransTxInfo web.TabtransInfoApx
		if er = rows.Scan(
			&tabtransTxInfo.Tgl_Trans,
			&tabtransTxInfo.Kode_LKM,
			&tabtransTxInfo.Kode_Trans,
			&tabtransTxInfo.Idpel,
			&tabtransTxInfo.Kuitansi,
			&tabtransTxInfo.Pokok,
			&tabtransTxInfo.Keterangan,
			&tabtransTxInfo.Product_Code,
			&tabtransTxInfo.Biller_Code,
			&tabtransTxInfo.User_Id,
		); er != nil {
			return listTx, er
		}

		listTx = append(listTx, tabtransTxInfo)
	}

	if len(listTx) == 0 {
		return listTx, err.NoRecord
	} else {
		return
	}
}

func (a *ApexTransRepoMysqlImpl) GetPrimaryTrxBelongToRecreateApx(kuitansi, MyKdTrans, BankCode string) (transApx entities.TransApx, er error) {
	row := a.apexDb.QueryRow(`SELECT 
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
		pay_lkm_source,
		pay_lkm_norek,
		pay_idpel,
		pay_biller_code,
		pay_product_code
	FROM tabtrans WHERE kuitansi= ? AND my_kode_trans= ? AND pay_lkm_source = ? LIMIT 1`, kuitansi, MyKdTrans, BankCode)
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
			return transApx, errors.New(fmt.Sprint("error while get lkm credit transfer data: ", er.Error()))
		}
	}
	return
}

func (a *ApexTransRepoMysqlImpl) DuplicateTrxBelongToRecreateApx(copy entities.TransApx) (er error) {

	apexTransRepo, _ := NewApexTransRepo()

	_, er = apexTransRepo.GetPrimaryTrxBelongToRecreateApx(copy.Kuitansi, "200", copy.No_rekening)
	if er != nil {
		stmt, er := a.apexDb.Prepare(`INSERT INTO tabtrans(
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
			pay_lkm_source,
			pay_lkm_norek,
			pay_idpel,
			pay_biller_code,
			pay_product_code
		) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`)
		if er != nil {
			return errors.New(fmt.Sprint("error while prepare add tabtrans transaction: ", er.Error()))
		}
		defer func() {
			_ = stmt.Close()
		}()

		// Get Trans ID
		tabtransId, err := apexTransRepo.GetTransIdApx()
		if er != nil {
			return err
		}

		if _, er := stmt.Exec(
			tabtransId,
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
			copy.Pay_lkm_source,
			copy.Pay_lkm_norek,
			copy.Pay_idpel,
			copy.Pay_biller_code,
			copy.Pay_product_code); er != nil {
			return errors.New(fmt.Sprint("error while add tabtrans transaction: ", er.Error()))
		}

	} else {
		entities.PrintError(fmt.Sprint(err.DuplicateEntry, ", Apex: ", copy.Kuitansi, "(290) is exist"))
		return err.DuplicateEntry

	}

	return nil
}

func (a *ApexTransRepoMysqlImpl) GetTabtransListByStanApx(stan string) (listTx []web.TabtransInfoApx, er error) {
	rows, er := a.apexDb.Query(`SELECT
		DATE_FORMAT(tgl_trans, "%d/%m/%Y") AS tgl_trans,
		no_rekening,
		kode_trans,
		pay_idpel,
		kuitansi,
		pokok,
		keterangan,
		pay_product_code,
		pay_biller_code,
		userid
	FROM tabtrans WHERE kuitansi LIKE "%` + stan + `%"`)
	if er != nil {
		return listTx, er
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var tabtransTxInfo web.TabtransInfoApx
		if er = rows.Scan(
			&tabtransTxInfo.Tgl_Trans,
			&tabtransTxInfo.Kode_LKM,
			&tabtransTxInfo.Kode_Trans,
			&tabtransTxInfo.Idpel,
			&tabtransTxInfo.Kuitansi,
			&tabtransTxInfo.Pokok,
			&tabtransTxInfo.Keterangan,
			&tabtransTxInfo.Product_Code,
			&tabtransTxInfo.Biller_Code,
			&tabtransTxInfo.User_Id,
		); er != nil {
			return listTx, er
		}

		listTx = append(listTx, tabtransTxInfo)
	}

	if len(listTx) == 0 {
		return listTx, err.NoRecord
	} else {
		return
	}
}

// ===================================================================================================================================================
// ==================================================================TEMPORARY FUNCTIONS==============================================================
// ===================================================================================================================================================
func (a *ApexTransRepoMysqlImpl) GetRekeningLKMByStatusActive() (lists []string, er error) {
	rows, er := a.apexDb.Query("SELECT no_rekening FROM tabung WHERE status = 1")
	if er != nil {
		return lists, er
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		var list entities.LKMlist
		if er = rows.Scan(&list.KodeLKM); er != nil {
			return lists, er
		}

		lists = append(lists, list.KodeLKM)
	}

	if len(lists) == 0 {
		return lists, err.NoRecord
	} else {
		return
	}
}

func (a *ApexTransRepoMysqlImpl) CalculateSaldoOnRekeningLKM(kodeLKM string) (data entities.CalculateSaldoResult, er error) {
	var tabtrans entities.RepostingData
	rows, err := a.apexDb.Query(`SELECT
	  tab.no_rekening,
	  SUM(CASE WHEN trans.my_kode_trans='100' THEN trans.pokok ELSE 0 END) AS total_kredit,
	  SUM(CASE WHEN trans.my_kode_trans='200' THEN trans.pokok ELSE 0 END) AS total_debet
	FROM tabung AS tab LEFT JOIN tabtrans AS trans ON (tab.no_rekening = trans.no_rekening)
	WHERE tab.no_rekening = ? GROUP BY tab.no_rekening`, kodeLKM)
	if err != nil {
		return data, er
	}
	for rows.Next() {
		rows.Scan(
			&tabtrans.KodeLKM,
			&tabtrans.TotalKredit,
			&tabtrans.TotalDebet,
		)
	}
	data.KodeLKM = tabtrans.KodeLKM
	data.SaldoAkhir = tabtrans.TotalKredit - tabtrans.TotalDebet
	return data, nil
}

func (a *ApexTransRepoMysqlImpl) RepostingSaldoOnRekeningLKMByScheduler(listOfKodeLKM ...string) (er error) {
	var wg sync.WaitGroup

	entities.PrintRepoChan <- entities.PrintRepo{Status: entities.PRINT_INIT_REPO_CHAN, Size: len(listOfKodeLKM)}

	for _, each := range listOfKodeLKM {

		wg.Add(1)

		go func(each string, w *sync.WaitGroup) {
			defer w.Done()

			var status = entities.PRINT_SUCCESS_STATUS_REPO_CHAN
			var msg = entities.PRINT_SUCCESS_MSG_REPO_CHAN

			er := a.doRepostingSaldoProcs(each)
			if er != nil {
				status = entities.PRINT_FAILED_STATUS_REPO_CHAN
				msg = er.Error()
			}
			var printRepo = entities.PrintRepo{
				KodeLKM: each,
				Status:  status,
				Message: msg,
			}
			entities.PrintRepoChan <- printRepo
		}(each, &wg)
	}
	wg.Wait()
	entities.PrintRepoChan <- entities.PrintRepo{Status: entities.PRINT_FINISH_REPO_CHAN}
	return
}

func (a *ApexTransRepoMysqlImpl) doRepostingSaldoProcs(data string) (er error) {
	lkm, er := a.CalculateSaldoOnRekeningLKM(data)
	if er != nil {
		return errors.New(fmt.Sprint("error while calculating saldo: ", er.Error()))
	}
	stmt, er := a.apexDb.Prepare(`UPDATE tabung SET saldo_akhir = ? WHERE no_rekening = ?`)
	if er != nil {
		return errors.New(fmt.Sprint("error while prepare reposting saldo: ", er.Error()))
	}
	defer func() {
		_ = stmt.Close()
	}()
	if _, er = stmt.Exec(
		lkm.SaldoAkhir,
		lkm.KodeLKM,
	); er != nil {
		return errors.New(fmt.Sprint("error while processing reposting saldo: ", er.Error()))
	}
	return
}
