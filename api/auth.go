package api

import (
	"bytes"
	"crypto/sha512"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Username string `json:"username" binging:"required"`
	Password string `json:"password" binging:"required"`
}
type loginResponse struct {
	UserID int32  `json:"user_id"`
	Token  string `json:"token"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (server *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
	}

	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = CompareHash(req.Password, user.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	expirationTime := time.Now().Add(time.Minute * 100)

	claims := &Claims{
		Username: req.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	generatedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := generatedToken.SignedString([]byte("secret"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := &loginResponse{
		UserID: user.ID,
		Token:  token,
	}

	ctx.JSON(http.StatusOK, arg)
}

func CompareHash(password string, hash string) (err error) {

	hashedInput := sha512.Sum512_256([]byte(password))
	trimmedHash := bytes.Trim(hashedInput[:], "\x00")
	preparedPassword := string(trimmedHash)
	plainTextInBytes := []byte(preparedPassword)
	hashTextInBytes := []byte(hash)
	err = bcrypt.CompareHashAndPassword(hashTextInBytes, plainTextInBytes)
	return
}
