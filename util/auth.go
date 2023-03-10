package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateToken(token string, ctx *gin.Context) {
	claims := &Claims{}
	tokenParse, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret_key"), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, err)
			return
		}
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	if !tokenParse.Valid {
		ctx.JSON(http.StatusUnauthorized, "Invalid token")
		return
	}

	ctx.Next()
}
