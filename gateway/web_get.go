package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"dev.sum7.eu/genofire/golang-lib/web"
)

type jsonDiscovery struct {
	UnifiedPush jsonUP `json:"unifiedpush"`
}
type jsonUP struct {
	Version int `json:"version"`
}

func Get(r *gin.Engine, ws *web.Service) {
	r.GET("/UP", func(c *gin.Context) {
		c.JSON(http.StatusOK, jsonDiscovery{
			jsonUP{
				Version: 1,
			},
		})
	})
}
