package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/web"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ChangeResponseCode(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)

	payload := web.ChangeResponseCode{}
	httpio.Bind(&payload)

	usecase := usecase.NewEchanneltransUsecase()
	er := usecase.ChangeResponseCode(payload)

	resp := web.RetransResponse{}
	if er != nil {
		entities.PrintError(er.Error())
		resp.ResponseCode = "1111"
		resp.ResponseMessage = er.Error()
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Change RC succeeded"
	}
	httpio.Response(http.StatusOK, resp)

}
