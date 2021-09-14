package main

import (
	"github.com/golang-jwt/jwt"
	"mellium.im/xmpp/jid"
)

// JWTSecret the secret
type JWTSecret string

// JWTToken data field
type JWTToken struct {
	jwt.StandardClaims
	PublicToken string `json:"token"`
	JID         string `json:"jid"`
}

// Generate an endpoint token by public token and jid
func (s JWTSecret) Generate(jid jid.JID, publicToken string) (string, error) {
	jwtToken := JWTToken{
		PublicToken: publicToken,
		JID:         jid.String(),
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtToken)
	endpointToken, err := claim.SignedString([]byte(s))
	if err != nil {
		return "", err
	}
	return endpointToken, nil
}

// Read endpoint token to public token and jid
func (s JWTSecret) Read(endpointToken string) (jid.JID, string, error) {
	token, err := jwt.ParseWithClaims(endpointToken, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s), nil
	})
	if err != nil {
		return jid.JID{}, "", err
	}
	claims, ok := token.Claims.(*JWTToken)
	if !ok {
		return jid.JID{}, "", jwt.ErrInvalidKey
	}
	addr, err := jid.Parse(claims.JID)
	if err != nil {
		return jid.JID{}, "", err
	}
	return addr, claims.PublicToken, nil
}
