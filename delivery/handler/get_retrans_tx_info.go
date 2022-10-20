package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/statuscode"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetRetransInfo(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.StanFilter{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	txInfo, er := usecase.GetRetransTxInfo(payload.Stan)
	if er != nil {
		if er == err.NoRecord {
			httpio.ResponseString(statuscode.StatusNoRecord, "record not found")
		} else {
			entities.PrintError(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
		}
	} else {
		httpio.Response(http.StatusOK, txInfo)
	}
}
