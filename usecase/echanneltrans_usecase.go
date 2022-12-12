package usecase

import (
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/repository/echanneltransrepo"
)

type EchanneltransUsecase interface {
	ChangeResponseCode(web.ChangeResponseCode) error
	UpdateIsoMsg(web.UpdateIsoMsg) error
}

type eChanneltransUsecase struct{}

func NewEchanneltransUsecase() EchanneltransUsecase {
	return &eChanneltransUsecase{}
}

func (e *eChanneltransUsecase) ChangeResponseCode(payload web.ChangeResponseCode) (er error) {

	// if (payload.Stan == "") || (payload.RC == "") {
	// 	return err.FieldMustBeExist
	// }

	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	er = dataRepo.ChangeResponseCode(payload.RC, payload.Stan, 0)
	if er != nil {
		return er
	}
	return nil
}

func (e *eChanneltransUsecase) UpdateIsoMsg(payload web.UpdateIsoMsg) (er error) {

	// if (payload.Iso_Msg == "") || (payload.Stan == "") {
	// 	return err.FieldMustBeExist
	// }

	dataRepo, _ := echanneltransrepo.NewEchannelTransRepo()
	er = dataRepo.UpdateIsoMsg(payload.Iso_Msg, payload.Stan)
	if er != nil {
		return er
	}
	return nil
}
