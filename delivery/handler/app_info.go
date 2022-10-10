package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func AppInfo(ctx *gin.Context) {

	httpio := httpio.NewRequestIO(ctx)
	httpio.Recv()

	appInfo := map[string]interface{}{
		"App Name":        os.Getenv("application.name"),
		"App Description": os.Getenv("application.desc"),
		"App Version":     os.Getenv("application.version"),
		"App Author":      os.Getenv("application.author"),
	}

	httpio.Response(http.StatusOK, appInfo)
}
