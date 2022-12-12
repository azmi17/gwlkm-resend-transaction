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

func ResendReversedTransByStan(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.StanFilter{}
	rerr := httpio.BindWithErr(&payload)
	if rerr != nil {
		errors := helper.FormatValidationError(rerr)
		errorMesage := gin.H{"errors": errors}
		response := helper.ApiResponse("Resend gwlkm transaction failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		httpio.Response(http.StatusUnprocessableEntity, response)
		return
	}

	usecase := usecase.NewRetransactionUsecase()
	newStan, er := usecase.ResendGwlkmTransaction(payload.Stan)

	resp := web.RetransWithNewStanResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
		resp.NewStan = newStan
		if resp.ResponseMessage == helper.AlreadyTransacted {
			resp.ResponseCode = "0000"
		}
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Resend gwlkm transaction succeeded"
		resp.NewStan = newStan
	}
	httpio.Response(http.StatusOK, resp)
}
