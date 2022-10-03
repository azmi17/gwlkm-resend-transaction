package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"

	"github.com/gin-gonic/gin"
)

func DeleteHandler(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx)
	data := entities.EmployeeFilter{}

	httpio.Bind(&data)

	ctx.JSON(200, data)
	httpio.Response(200, data)
}
