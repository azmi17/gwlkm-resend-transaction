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
	// err := httpio.BindWithErr(payload)
	// err := ctx.ShouldBind(&payload)
	// if err != nil {
	// 	errors := helper.FormatValidationError(err)
	// 	errorMesage := gin.H{"errors": errors}

	// 	response := helper.ApiResponse("Failed to update", http.StatusUnprocessableEntity, "failed", errorMesage)
	// 	ctx.JSON(http.StatusUnprocessableEntity, response)
	// 	return
	// }
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
