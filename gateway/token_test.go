package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"mellium.im/xmpp/jid"
)

func TestJWT(t *testing.T) {
	assert := assert.New(t)

	addr := "a@example.org"
	token := "pushtoken"

	secret := JWTSecret("CHANGEME")
	jwt, err := secret.Generate(jid.MustParse(addr), token)
	assert.NoError(err)
	assert.NotEqual("", jwt)

	jid, iToken, err := secret.Read(jwt)
	assert.NoError(err)
	assert.Equal(addr, jid.String())
	assert.Equal(iToken, token)
}
