package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/entities"
	"gwlkm-resend-transaction/entities/err"
	"gwlkm-resend-transaction/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Temporary Functions
func RepostingAllByApi(ctx *gin.Context) {

	httpio := httpio.NewRequestIO(ctx)
	httpio.Recv()

	usecase := usecase.NewApexTransUsecase()
	er := usecase.RepostingSaldoApexByScheduler()

	// go RepostingSchedulerRepoObserver()

	resp := entities.SchedulerResponse{}
	if er != nil {
		if er == err.NoRecord {
			resp.ResponseCode = "1111"
			resp.ResponseMessage = er.Error()
		} else {
			entities.PrintError(er.Error())
			httpio.ResponseString(http.StatusInternalServerError, "internal service error")
			return
		}
	} else {
		resp.ResponseCode = "0000"
		resp.ResponseMessage = "Reposting saldo succeeded"
	}

	httpio.Response(http.StatusOK, resp)
}
