package middleware

import (
	"net/http"
	"os"
	"time"

	"github.com/eFournierRobert/themedia/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Authorization(context *gin.Context) {
	tokenString, err := context.Cookie("Authorization")
	if err != nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || token == nil {
		context.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if claims != nil && ok {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if !tools.DoesUserExist(claims["sub"].(string)) {
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		context.Set("userUUID", claims["sub"].(string))
		context.Next()
	} else {
		context.AbortWithStatus(http.StatusUnauthorized)
	}
}
