package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecycleSuspectRevBiller(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.SuspectRevBiller{}
	rerr := httpio.BindWithErr(&payload)
	if rerr != nil {
		errors := helper.FormatValidationError(rerr)
		errorMesage := gin.H{"errors": errors}
		response := helper.ApiResponse("Execute suspect reversal biller failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		httpio.Response(http.StatusUnprocessableEntity, response)
		return
	}

	usecase := usecase.NewRetransactionUsecase()
	er := usecase.RecycleSuspectRevBiller(payload)

	resp := web.RetransResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseMessage = er.Error()
		resp.ResponseCode = "1111"
		if resp.ResponseMessage == helper.AlreadyReversed {
			resp.ResponseCode = "0000"
		}
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Execute suspect reversal biller succeeded"
	}
	httpio.Response(http.StatusOK, resp)
}
