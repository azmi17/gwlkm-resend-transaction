package apextransrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
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

func (a *ApexTransRepoMysqlImpl) GetTxInfoApx(kuitansi string) (transApx entities.TransApx, er error) {
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
			item.Tgl_real_trans,
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
			item.Tgl_real_trans,
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

func (a *ApexTransRepoMysqlImpl) DeleteTxApx(kuitansi string) (er error) {

	stmt, er := a.apexDb.Prepare("DELETE FROM tabtrans WHERE kuitansi = ?")
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