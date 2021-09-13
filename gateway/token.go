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
	Token string `json:"token"`
	JID   string `json:"jid"`
}

// Generate an jwt token by token and jid
func (s JWTSecret) Generate(jid jid.JID, token string) (string, error) {
	jwtToken := JWTToken{
		Token: token,
		JID:   jid.String(),
	}
	claim := jwt.NewWithClaims(jwt.SigningMethodHS512, jwtToken)
	t, err := claim.SignedString([]byte(s))
	if err != nil {
		return "", err
	}
	return t, nil
}

// Read token to token and jid
func (s JWTSecret) Read(jwtToken string) (jid.JID, string, error) {
	token, err := jwt.ParseWithClaims(jwtToken, &JWTToken{}, func(token *jwt.Token) (interface{}, error) {
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
	return addr, claims.Token, nil
}
