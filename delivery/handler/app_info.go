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
		"App Name":         "e-Channel Recycle Transaction",
		"App Description":  "e-Channel Retransaction API Endpoint",
		"App Version":      "1.4.2",
		"App Author":       "Azmi Farhan",
		"App Release Date": "13/10/2022 17:53",
	}

	httpio.Response(http.StatusOK, appInfo)
}
