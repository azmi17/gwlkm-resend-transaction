package handler

import (
	"gwlkm-resend-transaction/delivery/handler/httpio"
	"gwlkm-resend-transaction/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppInfo(ctx *gin.Context) {

	httpio := httpio.NewRequestIO(ctx)
	httpio.Recv()

	appInfo := map[string]interface{}{
		"App Name":        helper.AppName,
		"App Description": helper.AppDescription,
		"App Version":     helper.AppVersion,
		"App Author":      helper.AppAuthor,
		"App Updated At:": helper.LastBuild,
	}

	httpio.Response(http.StatusOK, appInfo)
}
