package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResendReversedTransByStan(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.StanFilter{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	newStan, er := usecase.ResendGwlkmTransaction(payload.Stan)

	resp := web.RetransWithNewStanResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
		resp.NewStan = newStan
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Resend gwlkm transaction succeeded"
		resp.NewStan = newStan
	}
	httpio.Response(http.StatusOK, resp)
}
