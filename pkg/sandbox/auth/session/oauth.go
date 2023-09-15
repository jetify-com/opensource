package session

import (
	"golang.org/x/oauth2"
)

func FromOauth2(otok *oauth2.Token) (*Token, error) {
	if otok == nil {
		return nil, nil
	}

	tok := Token{
		Token:   *otok,
		IDToken: getRawIDToken(otok),
	}

	claims := tok.IDClaims()
	if claims != nil {
		tok.Issuer = claims.Issuer
	}

	return &tok, nil
}

func getRawIDToken(otok *oauth2.Token) string {
	if otok == nil {
		return ""
	}

	rawIDTok := otok.Extra("id_token")
	if rawIDTok == nil {
		return ""
	}
	return otok.Extra("id_token").(string)
}
