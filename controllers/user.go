package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gimtwi/go-jwt-auth/initializers"
	"github.com/gimtwi/go-jwt-auth/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(ctx *gin.Context) {

	if ctx.Bind(&types.UserRequest) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(types.UserRequest.Password), 10)
	fmt.Println(hash)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to hash the password",
		})
		return
	}

	user := types.User{Email: types.UserRequest.Email, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create user",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})

}

func Login(ctx *gin.Context) {

	if ctx.Bind(&types.UserRequest) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to read body",
		})
		return
	}

	var user types.User

	initializers.DB.First(&user, "email = ?", types.UserRequest.Email)

	if user.ID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	errPWD := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(types.UserRequest.Password))

	if errPWD != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(), //30 days
	})
	secret := os.Getenv("JWT_SECRET")

	tokenStr, err := token.SignedString([]byte(secret))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to create token",
		})
		return
	}

	ctx.SetSameSite(http.SameSiteLaxMode)

	ctx.SetCookie("auth", tokenStr, 3600*24*30, "", "", false, true) //secure false only on localhost, change to true in prod

	ctx.JSON(http.StatusOK, gin.H{})

}

func Validate(ctx *gin.Context) {
	user, _ := ctx.Get("user")

	ctx.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
