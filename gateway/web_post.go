package main

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/web"
)

func Post(r *gin.Engine, ws *web.Service, xmpp *XMPPService, jwtsecret JWTSecret) {
	r.POST("/UP", func(c *gin.Context) {
		to, token, err := jwtsecret.Read(c.Query("token"))
		if err != nil {
			c.JSON(http.StatusUnauthorized, web.HTTPError{
				Message: "jwt token unauthoried - or not given",
				Error:   err.Error(),
			})
			return
		}
		b, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, web.HTTPError{
				Message: web.ErrAPIInvalidRequestFormat.Error(),
				Error:   err.Error(),
			})
			return
		}
		content := string(b)
		if err := xmpp.SendMessage(to, token, content); err != nil {
			c.JSON(http.StatusNotFound, web.HTTPError{
				Message: "unable to forward to xmpp",
				Error:   err.Error(),
			})
			return
		}
		c.JSON(http.StatusAccepted, content)
	})
}
