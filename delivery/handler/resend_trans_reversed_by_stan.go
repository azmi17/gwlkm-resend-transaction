package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResendReversedTransByStan(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := entities.TransHistoryRequest{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	newStan, er := usecase.ResendReversedTransaction(payload.Stan)

	resp := entities.TransHistoryResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
		resp.NewStan = newStan
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Resend transaction succeeded"
		resp.NewStan = newStan
	}
	httpio.Response(http.StatusOK, resp)
}
