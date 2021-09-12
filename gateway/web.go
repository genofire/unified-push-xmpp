package main

import (
	"github.com/gin-gonic/gin"
	"dev.sum7.eu/genofire/golang-lib/web"
	"dev.sum7.eu/genofire/golang-lib/web/api/status"
	"dev.sum7.eu/genofire/golang-lib/web/metrics"
)

// Bind to webservice
// @title UnifiedPush API for XMPP
// @version 1.0
// @description This is the first version of an UnifiedPush Gateway for XMPP
// @termsOfService http://swagger.io/terms/
// -host up.chat.sum7.eu
// @BasePath /
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func Bind(xmpp *XMPPService) web.ModuleRegisterFunc {
	return func(r *gin.Engine, ws *web.Service) {
		// docs.Bind(r, ws)

		status.Register(r, ws)
		metrics.Register(r, ws)
		Get(r, ws)
		Post(r, ws, xmpp)
	}
}
