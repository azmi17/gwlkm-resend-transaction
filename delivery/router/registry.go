/*
 * Copyright (c) 2022 Randy Ardiansyah https://github.com/randyardiansyah25/<repo>
 *
 * Created Date: Wednesday, 16/03/2022, 10:32:08
 * Author: Randy Ardiansyah
 *
 * Filename: /home/Documents/workspace/go/src/router-template/delivery/router/registry.go
 * Project : /home/Documents/workspace/go/src/router-template/delivery/router
 *
 * HISTORY:
 * Date                  	By                 	Comments
 * ----------------------	-------------------	--------------------------------------------------------------------------------------------------------------------
 */

package router

import (
	"gwlkm-resend-transaction/delivery/handler"

	"github.com/gin-gonic/gin"
)

func RegisterHandler(router *gin.Engine) {

	// API Versioning:
	apiv1 := router.Group("api/v1/retransaction")

	// API Endpoint:
	apiv1.GET("/version", handler.AppInfo)
	apiv1.POST("/resend", handler.ResendReversedTransByStan)
	apiv1.POST("/history", handler.GetRetransInfo)
	apiv1.PUT("/change-rc", handler.ChangeResponseCode)
	apiv1.PUT("/isomsg", handler.UpdateIsoMsg)
	apiv1.POST("/tabtrans", handler.GetTabtransInfo)
	apiv1.POST("/tabtrans-bystan", handler.GetTabtransInfoByStan)

	apiv1.POST("/resent", handler.ResendTransByStan)
	apiv1.POST("/tintcr/reversal", handler.ResendLkmTransferSMprematureRevOnCre)
	apiv1.PUT("/reset-password-apex", handler.ResetApexPassword)

	apiv1.POST("/resend/reversal", handler.ResendReversalGwlkm)
	apiv1.POST("/apex/recreate", handler.RecreateSuccessTransactionApx)

	// Temporary functions
	apiv1.POST("/repostings/all", handler.RepostingAllByApi)
}
