package controllers

import (
	"net/http"
	"os"
	"strconv"

	"errors"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func first[T any](val T, _ error) T {
	return val
}

func Signup(ctx *gin.Context) {
	var req SignupRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		first(strconv.Atoi(os.Getenv("BCRYPT_COST"))),
	)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			badRequest(ctx, "PasswordTooLong")
			return
		}
		unknownError(ctx)
		return
	}

	user := models.User{Email: req.Email, PasswordHash: hash}
	err = db.Create(&user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			badRequest(ctx, "EmailAlreadyInUse")
			return
		}
		unknownError(ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "Success"})
}

func Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
		return
	}

	var user models.User
	db.FindOneWhere(&models.User{Email: req.Email}, &user)
	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(req.Password)); err != nil {
		badRequest(ctx, "InvalidCredentials")
		return
	}

	signedAccessToken, signedRefreshToken := createTokenPair(user.ID)
	ctx.JSON(http.StatusOK, LoginResponse{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken})

}
