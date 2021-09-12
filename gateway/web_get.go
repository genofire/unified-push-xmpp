package main

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"dev.sum7.eu/genofire/golang-lib/web"
)

type jsonDiscovery struct {
	Unifiedpush jsonUP `json:"unifiedpush"`
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
