package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gimtwi/go-jwt-auth/initializers"
	"github.com/gimtwi/go-jwt-auth/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(ctx *gin.Context) {
	tokenStr, err := ctx.Cookie("auth")

	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	secret := os.Getenv("JWT_SECRET")

	token, _ := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		var user types.User

		initializers.DB.First(&user, claims["sub"])
		if user.ID == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		ctx.Set("user", user)

		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	ctx.Next()
}
