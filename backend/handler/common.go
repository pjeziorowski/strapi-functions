package handler

import "github.com/dgrijalva/jwt-go"

type JwtClaims struct {
	HasuraClaims HasuraClaims `json:"https://hasura.io/jwt/claims"`
	jwt.StandardClaims
}

type HasuraClaims struct {
	XHasuraUserId       string   `json:"x-hasura-user-id"`
	XHasuraDefaultRole  string   `json:"x-hasura-default-role"`
	XHasuraAllowedRoles []string `json:"x-hasura-allowed-roles"`
}
