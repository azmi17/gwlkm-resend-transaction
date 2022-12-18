package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecreateReversalTransactionApx(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.RecreateApexRequest{}
	rerr := httpio.BindWithErr(&payload)
	if rerr != nil {
		errors := helper.FormatValidationError(rerr)
		errorMesage := gin.H{"errors": errors}
		response := helper.ApiResponse("Recreate apex reversal transaction failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		httpio.Response(http.StatusUnprocessableEntity, response)
		return
	}

	usecase := usecase.NewApexTransUsecase()
	er := usecase.RecreateReversalTransactionApx(payload)

	resp := web.RetransResponse{}
	if er != nil {
		if er == err.NoRecord || er == err.DuplicateEntry {
			resp.ResponseCode = "1111"
			resp.ResponseMessage = er.Error()
		} else {
			entities.PrintError(er.Error())
			entities.PrintLog(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
			return
		}
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Recreate apex reversal transaction succeeded"
	}

	httpio.Response(http.StatusOK, resp)
}
