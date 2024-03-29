package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"errors"

	"github.com/AnshVM/flashpoll-backend/db"
	"github.com/AnshVM/flashpoll-backend/models"
	"github.com/AnshVM/flashpoll-backend/types"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RefreshRequest = types.RefreshRequest
type SignupRequest = types.SignupRequest
type LoginRequest = types.LoginRequest
type LoginResponse = types.LoginResponse

func first[T any](val T, _ error) T {
	return val
}

var maxCookieAge = 24 * 60 * 60

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
			badRequest(ctx, "Password too long")
			return
		}
		unknownError(ctx)
		return
	}

	user := models.User{Email: req.Email, PasswordHash: hash, Username: req.Username}
	err = db.Create(&user)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			var dupUserWithEmail models.User
			result := db.FindOneWhere(&models.User{Email: req.Email}, &dupUserWithEmail)
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				badRequest(ctx, "DuplicateUsername")
				return
			}
			badRequest(ctx, "DuplicateEmail")
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
	ctx.SetCookie("accessToken", signedAccessToken, maxCookieAge, "/", os.Getenv("SERVER_DOMAIN"), true, true)
	ctx.SetCookie("refreshToken", signedRefreshToken, maxCookieAge*10, "/", os.Getenv("SERVER_DOMAIN"), true, true)

	ctx.JSON(http.StatusOK, LoginResponse{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken})
}

func RefreshTokens(ctx *gin.Context) {

	fmt.Println("hererer")
	// refreshToken, err := ctx.Cookie("refreshToken")

	// if err != nil {
	// 	unauthorized(ctx)
	// 	return
	// }

	var req RefreshRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		invalidRequestBody(ctx)
		return
	}

	claims, err := parseToken(req.RefreshToken, []byte(os.Getenv("REFRESH_TOKENS_SECRET_KEY")))

	fmt.Printf("%v", claims)

	if err != nil {
		unauthorized(ctx)
		return
	}

	signedAccessToken, signedRefreshToken := createTokenPair(claims.UserID)

	ctx.SetCookie("accessToken", signedAccessToken, maxCookieAge, "/", os.Getenv("SERVER_DOMAIN"), true, true)
	ctx.SetCookie("refreshToken", signedRefreshToken, maxCookieAge*10, "/", os.Getenv("SERVER_DOMAIN"), true, true)
	fmt.Println("accessToken")
	fmt.Println(signedAccessToken)
	ctx.JSON(http.StatusOK, LoginResponse{AccessToken: signedAccessToken, RefreshToken: signedRefreshToken})
}

func Logout(ctx *gin.Context) {
	ctx.SetCookie("accessToken", "", -1, "/", os.Getenv("SERVER_DOMAIN"), true, true)
	ctx.SetCookie("refreshToken", "", -1, "/", os.Getenv("SERVER_DOMAIN"), true, true)
	ctx.JSON(http.StatusOK, "USER_LOGGED_OUT")
}
