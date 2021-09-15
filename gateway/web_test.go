package main

import (
	"net/http"
	"testing"

	"dev.sum7.eu/genofire/golang-lib/web/webtest"
	"github.com/stretchr/testify/assert"
)

func TestWebGet(t *testing.T) {
	assert := assert.New(t)
	s, err := webtest.New(Bind(nil, ""))
	assert.NoError(err)
	defer s.Close()
	assert.NotNil(s)

	up := jsonDiscovery{}
	// GET
	err = s.Request(http.MethodGet, "/UP", nil, http.StatusOK, &up)
	assert.NoError(err)
	assert.Equal(1, up.UnifiedPush.Version)
}
