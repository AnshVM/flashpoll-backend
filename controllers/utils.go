package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func badRequest(ctx *gin.Context, err string) {
	ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
}

func invalidRequestBody(ctx *gin.Context) {
	badRequest(ctx, "InvalidRequestBody")
}

func unknownError(ctx *gin.Context) {
	badRequest(ctx, "UnknownError")
}

type Claims struct {
	UserID uint `json:"userID"`
	jwt.RegisteredClaims
}

func createTokenPair(id uint) (string, string) {
	accessTokenclaims := Claims{
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	refreshTokenClaims := Claims{
		id,
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	signedAccessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims).SignedString([]byte(os.Getenv("ACCESS_TOKENS_SECRET_KEY")))
	if err != nil {
		panic(err)
	}
	signedRefreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString([]byte(os.Getenv("REFRESH_TOKENS_SECRET_KEY")))
	if err != nil {
		panic(err)
	}
	return signedAccessToken, signedRefreshToken
}
