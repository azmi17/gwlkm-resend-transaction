package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResendTransByStan(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.StanFilter{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	er := usecase.ResendTransaction(payload.Stan)

	resp := web.RetransResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Resend transaction succeeded"
	}
	httpio.Response(http.StatusOK, resp)
}
