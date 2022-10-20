package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UpdateIsoMsg(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.UpdateIsoMsg{}
	httpio.Bind(&payload)

	usecase := usecase.NewRetransactionUsecase()
	er := usecase.UpdateIsoMsg(payload)

	resp := web.RetransResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Update iso message succeeded"
	}
	httpio.Response(http.StatusOK, resp)

}
