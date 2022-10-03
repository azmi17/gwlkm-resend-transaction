package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/entities/statuscode"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetEmployee(ctx *gin.Context) {
	httpio := httpio.NewRequestIO(ctx) // catch request
	httpio.Recv()

	payload := entities.EmployeeFilter{}
	httpio.BindJSON(&payload)

	ucase := usecase.NewEmployeeUsecase()
	employee, er := ucase.GetEmployee(int(payload.Id))
	if er != nil {
		if er == err.DuplicateEntry {
			httpio.ResponseString(statuscode.StatusDuplicate, "Data karyawan sudah tersedia!")
		} else {
			entities.PrintError(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
		}
	} else {
		httpio.Response(http.StatusOK, employee)
	}
}
