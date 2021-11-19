package auth

import (
	"fmt"
	"functions/backend/config"
	"functions/backend/handler"
	"github.com/dgrijalva/jwt-go"
	"github.com/hasura/go-graphql-client"
	"time"
)

func SignTokenFor(userId graphql.Int) (string, error) {
	claims := &handler.JwtClaims{
		handler.HasuraClaims{
			XHasuraUserId:       fmt.Sprintf("%v", userId),
			XHasuraDefaultRole:  "admin",
			XHasuraAllowedRoles: []string{"admin"},
		},
		jwt.StandardClaims{ExpiresAt: time.Now().Add(time.Hour * 72).Unix()},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(config.Secret))

	if err != nil {
		return "", err
	}

	return signed, nil
}
