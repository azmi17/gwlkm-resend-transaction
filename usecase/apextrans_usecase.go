package usecase

import (
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/repository/apextransrepo"
	"gwlkm-resend-transaction/repository/echanneltransrepo"
	"time"

	"github.com/kpango/glg"
)

type ApexTransUsecase interface {
	GetTabtransListApx(kuitansi string) ([]web.TabtransInfoApx, error)
	GetTabtransListByStanApx(stan string) ([]web.TabtransInfoApx, error)
	RecreateSuccessTransactionApx(request web.RecreateApexRequest) error
	RecreateReversalTransactionApx(request web.RecreateApexRequest) error

	RepostingSaldoApexByScheduler() (er error) // Temporary Functions
}

type apextransUsecase struct{}

func NewApexTransUsecase() ApexTransUsecase {
	return &apextransUsecase{}
}

func (a *apextransUsecase) GetTabtransListApx(kuitansi string) (detailTx []web.TabtransInfoApx, er error) {
	repo, _ := apextransrepo.NewApexTransRepo()

	detailTx, er = repo.GetTabtransListApx(kuitansi)
	if er != nil {
		return detailTx, er
	}

	if len(detailTx) == 0 {
		return detailTx, err.NoRecord
	}

	return detailTx, nil
}

func (a *apextransUsecase) GetTabtransListByStanApx(stan string) (detailTx []web.TabtransInfoApx, er error) {
	repo, _ := apextransrepo.NewApexTransRepo()

	detailTx, er = repo.GetTabtransListByStanApx(stan)
	if er != nil {
		return detailTx, er
	}

	if len(detailTx) == 0 {
		return detailTx, err.NoRecord
	}

	return detailTx, nil
}

func (a *apextransUsecase) RecreateSuccessTransactionApx(request web.RecreateApexRequest) (er error) {
	apxRepo, _ := apextransrepo.NewApexTransRepo()
	echRepo, _ := echanneltransrepo.NewEchannelTransRepo()

	trxSource, er := echRepo.GetOriginData(request.Stan)
	if er != nil {
		return er
	}

	trans := entities.TransApx{}
	trans.Tgl_trans, _ = time.Parse("2006-01-02", trxSource.Tgl_Trans_Str[0:4]+"-"+trxSource.Tgl_Trans_Str[4:6]+"-"+trxSource.Tgl_Trans_Str[6:8])
	trans.No_rekening = request.KodeLKM
	trans.Kode_trans = "290"
	trans.My_kode_trans = "200"
	trans.Pokok = float64(trxSource.Amount + trxSource.Profit_Excluded + trxSource.Profit_Included)
	trans.Kuitansi = trxSource.Ref
	trans.Userid = 121
	trans.Keterangan = "Echannel, norek: " + trxSource.Rek_Id + ", (" + trxSource.Biller_Code + "-" + trxSource.Product_Code + "), Idpel " + trxSource.Subscriber_Id
	trans.Verifikasi = "1"
	trans.Tob = "T"
	trans.Sandi_trans = "PAY"
	trans.Posted_to_gl = "0"
	trans.Kode_kantor = "001"
	trans.Jam = "08:00:00"
	trans.Pay_lkm_source = trxSource.Bank_Code
	trans.Pay_lkm_norek = trxSource.Rek_Id
	trans.Pay_idpel = trxSource.Subscriber_Id
	trans.Pay_biller_code = trxSource.Biller_Code
	trans.Pay_product_code = trxSource.Product_Code

	if er = apxRepo.DuplicateTrxBelongToRecreateApx(trans); er != nil {
		return er
	}

	return
}

func (a *apextransUsecase) RecreateReversalTransactionApx(request web.RecreateApexRequest) (er error) {
	apxRepo, _ := apextransrepo.NewApexTransRepo()
	echRepo, _ := echanneltransrepo.NewEchannelTransRepo()

	trxSource, er := echRepo.GetOriginData(request.Stan)
	if er != nil {
		return er
	}

	trans := entities.TransApx{}
	trans.Tgl_trans, _ = time.Parse("2006-01-02", trxSource.Tgl_Trans_Str[0:4]+"-"+trxSource.Tgl_Trans_Str[4:6]+"-"+trxSource.Tgl_Trans_Str[6:8])
	trans.No_rekening = request.KodeLKM
	trans.Kode_trans = "190"
	trans.My_kode_trans = "100"
	trans.Pokok = float64(trxSource.Amount + trxSource.Profit_Excluded + trxSource.Profit_Included)
	trans.Kuitansi = trxSource.Ref
	trans.Userid = 121
	trans.Keterangan = "Reversal Echannel, norek: " + trxSource.Rek_Id + ", (" + trxSource.Biller_Code + "-" + trxSource.Product_Code + "), Idpel " + trxSource.Subscriber_Id
	trans.Verifikasi = "1"
	trans.Tob = "T"
	trans.Sandi_trans = "PAY"
	trans.Posted_to_gl = "0"
	trans.Kode_kantor = "001"
	trans.Jam = "08:00:00"
	trans.Pay_lkm_source = trxSource.Bank_Code
	trans.Pay_lkm_norek = trxSource.Rek_Id
	trans.Pay_idpel = trxSource.Subscriber_Id
	trans.Pay_biller_code = trxSource.Biller_Code
	trans.Pay_product_code = trxSource.Product_Code

	if er = apxRepo.DuplicateTrxBelongToRecreateApx(trans); er != nil {
		return er
	}

	if er = echRepo.ChangeResponseCode("1100", request.Stan, trxSource.Trans_id); er != nil {
		return er
	}

	return
}

// Temporary Functions
func (a *apextransUsecase) RepostingSaldoApexByScheduler() (er error) {
	repo, _ := apextransrepo.NewApexTransRepo()

	list, er := repo.GetRekeningLKMByStatusActive()
	if er != nil {
		return er
	}

	_ = glg.Log("Reposting saldo apex is processing..")
	er = repo.RepostingSaldoOnRekeningLKMByScheduler(list...)
	if er != nil {
		return er
	}

	return nil
}
