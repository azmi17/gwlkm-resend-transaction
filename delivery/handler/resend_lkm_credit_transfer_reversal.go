package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResendLkmTransferSMprematureRevOnCre(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.StanFilter{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	er := usecase.ResendLkmTransferSMprematureRevOnCre(payload.Stan)

	resp := web.RetransResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseMessage = er.Error()
		resp.ResponseCode = "1111"
		if resp.ResponseMessage == "44-Transaksi Sudah di reversal!" {
			resp.ResponseCode = "0000"
		}

	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Resend TINTCR reversal succeeded"
	}
	httpio.Response(http.StatusOK, resp)
}
