package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/helper"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResetApexPassword(ctx *gin.Context) {

	// Init HTTP Request..
	httpio := httpio.NewRequestIO(ctx)

	// Call Payload and binding form
	payload := entities.KodeLKMFilter{}
	rerr := httpio.BindWithErr(&payload)
	if rerr != nil {
		errors := helper.FormatValidationError(rerr)
		errorMesage := gin.H{"errors": errors}
		response := helper.ApiResponse("Reset apex user web password failed", http.StatusUnprocessableEntity, "failed", errorMesage)
		httpio.Response(http.StatusUnprocessableEntity, response)
		return
	}

	usecase := usecase.NewSysUserUsecase()
	lkmPwd, er := usecase.ResetSysUserPassword(payload)

	resp := entities.ResetPasswrodResponse{}
	if er != nil {
		if er == err.NoRecord || er == err.FieldMustBeExist {
			resp.Response_Code = "1111"
			resp.Response_Msg = er.Error()
		} else {
			entities.PrintError(er.Error())
			entities.PrintLog(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
			return
		}
	} else {
		resp.Response_Code = "0000"
		resp.Response_Msg = "Reset apex user web password succeeded"
		resp.Data = &lkmPwd
	}

	httpio.Response(http.StatusOK, resp)
}
