package util

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func validateToken(token string, ctx *gin.Context) error {
	claims := &Claims{}
	tokenParse, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, err)
			return err
		}
		ctx.JSON(http.StatusBadRequest, err)
		return err
	}

	if !tokenParse.Valid {
		ctx.JSON(http.StatusUnauthorized, "Invalid token")
		return err
	}

	ctx.Next()
	return nil
}

func GetTokenInHeaderAndVerify(ctx *gin.Context) error {
	authorizationHeaderKey := ctx.GetHeader("authorization")
	fields := strings.Fields(authorizationHeaderKey)
	tokenToValidate := fields[1]
	err := validateToken(tokenToValidate, ctx)
	if err != nil {
		return err
	}
	return nil
}
