package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppInfo(ctx *gin.Context) {

	httpio := httpio.NewRequestIO(ctx)
	httpio.Recv()

	appInfo := map[string]interface{}{
		"App Name":        "e-Channel Recycle Transaction",
		"App Description": "e-Channel Retransaction API Endpoint",
		"App Version":     "1.4.3",
		"App Author":      "Azmi Farhan",
		"App Updated At:": "14/10/2022 15:27",
	}

	httpio.Response(http.StatusOK, appInfo)
}
